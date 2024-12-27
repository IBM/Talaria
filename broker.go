package main

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"net/url"
	"opentalaria/config"
	"opentalaria/utils"
	"strconv"
	"strings"
)

var (
	// by default in KRaft mode, generated broker IDs start from reserved.broker.max.id + 1,
	// where reserved.broker.max.id=1000 if the property is not set.
	// KRaft mode is the default Kafka mode, since Kafka v3.3.1, so OpenTalaria will implement default settings in KRaft mode.
	RESERVED_BROKER_MAX_ID = 1000
)

// NewBroker returns a new instance of Broker.
// For now OpenTalaria does not support rack awareness, but this will change in the future.
func NewBroker() (config.Broker, error) {
	broker := config.Broker{}

	advListenerStr, _ := utils.GetEnvVar("advertised.listeners", "")
	advertisedListeners := strings.Split(strings.ReplaceAll(advListenerStr, " ", ""), ",")

	listenerStr, ok := utils.GetEnvVar("listeners", "")
	if !ok {
		return config.Broker{}, errors.New("no listeners set")
	}
	listeners := strings.Split(strings.ReplaceAll(listenerStr, " ", ""), ",")

	if len(advertisedListeners) == 0 {
		advertisedListeners = listeners
	}

	listenersArray, err := parseListeners(listeners)
	if err != nil {
		return config.Broker{}, err
	}
	broker.Listeners = append(broker.Listeners, listenersArray...)

	err = validateListeners(&broker)
	if err != nil {
		return config.Broker{}, err
	}

	advertisedListenersArr, err := parseListeners(advertisedListeners)
	if err != nil {
		return config.Broker{}, err
	}
	broker.AdvertisedListeners = append(broker.AdvertisedListeners, advertisedListenersArr...)

	err = validateAdvertisedListeners(&broker)
	if err != nil {
		return config.Broker{}, err
	}

	brokerIdSetting, _ := utils.GetEnvVar("broker.id", "-1")

	brokerId, err := strconv.Atoi(brokerIdSetting)
	if err != nil {
		return broker, fmt.Errorf("error parsing broker.id: %s", err)
	}

	// validate Broker ID
	if brokerId > RESERVED_BROKER_MAX_ID {
		return broker, fmt.Errorf("the configured node ID is greater than `reserved.broker.max.id`. Please adjust the `reserved.broker.max.id` setting. [%d > %d]",
			brokerId,
			RESERVED_BROKER_MAX_ID)
	}

	if brokerId == -1 {
		brokerId = RESERVED_BROKER_MAX_ID + 1
	}

	broker.BrokerID = int32(brokerId)

	return broker, nil
}

func parseListeners(listeners []string) ([]config.Listener, error) {
	result := []config.Listener{}

	for _, l := range listeners {
		if l == "" {
			continue
		}

		listener, err := parseListener(l)
		if err != nil {
			return []config.Listener{}, err
		}

		result = append(result, listener)
	}

	return result, nil
}

func parseListener(l string) (config.Listener, error) {
	listener, err := url.Parse(l)
	if err != nil {
		return config.Listener{}, err
	}

	// parse the security protocol from the url scheme.
	// If the protocol is unknown treat the scheme as broker name and check the listener.security.protocol.map
	listenerName, securityProtocol, err := getBrokerNameComponents(listener.Scheme)
	if err != nil {
		return config.Listener{}, err
	}

	host, port, err := net.SplitHostPort(listener.Host)
	if err != nil {
		return config.Listener{}, err
	}

	parsedPort, err := strconv.Atoi(port)
	if err != nil {
		return config.Listener{}, err
	}

	return config.Listener{
		Host:             host,
		Port:             int32(parsedPort),
		SecurityProtocol: securityProtocol,
		ListenerName:     listenerName,
	}, nil
}

// getBrokerNameComponents checks if the broker name, inferred from the URL schema is a valid security protocol.
// If not, it checks the listener.security.protocol.map for mapping for custom broker names and returns the broker name/security protocol pair.
// If no mapping is found in the case of custom broker name, the function returns an error.
func getBrokerNameComponents(s string) (string, config.SecurityProtocol, error) {
	securityProtocol, ok := config.ParseSecurityProtocol(s)

	if ok {
		return s, securityProtocol, nil
	} else {
		// the listener schema is not a known security protocol, treat is as broker name
		// and extract the security protocol from listener.security.protocol.map
		listenerSpmStr, _ := utils.GetEnvVar("listener.security.protocol.map", "")
		spm := strings.Split(strings.ReplaceAll(listenerSpmStr, " ", ""), ",")

		for _, sp := range spm {
			components := strings.Split(sp, ":")

			if strings.EqualFold(s, components[0]) {
				securityProtocol, ok := config.ParseSecurityProtocol(components[1])
				if !ok {
					return "", config.UNDEFINED_SECURITY_PROTOCOL, fmt.Errorf("unknown security protocol for listener %s", components[0])
				}

				return s, securityProtocol, nil
			}
		}
	}

	return "", config.UNDEFINED_SECURITY_PROTOCOL, fmt.Errorf("broker %s not found in listener.security.protocol.map", s)
}

// validateListeners performs common checks on the listeners as per Kafka specification https://kafka.apache.org/documentation/#brokerconfigs_listeners.
// Broker name and port have to be unique. The exception is if the host for two entries is IPv4 and IPv6 respectively.
func validateListeners(b *config.Broker) error {
	ports := map[int32]string{}
	listenerNames := map[string]string{}

	for _, listener := range b.Listeners {
		// Check uniqueness for ports
		if val, ok := ports[listener.Port]; ok {
			if areIpProtocolsSame(listener.Host, val) {
				return fmt.Errorf("listener port is not unique for listener %s", listener.ListenerName)
			}
		}

		// Check uniqueness for broker names
		if val, ok := listenerNames[listener.ListenerName]; ok {
			if areIpProtocolsSame(listener.Host, val) {
				return fmt.Errorf("listener name is not unique for listener %s", listener.ListenerName)
			}
		}

		ports[listener.Port] = listener.Host
		listenerNames[listener.ListenerName] = listener.Host
	}

	return nil
}

func areIpProtocolsSame(host1, host2 string) bool {
	// ignore errors from ParseAddr, which will be thrown if a hostname is provided, we care only about IP addresses.
	addr1, _ := netip.ParseAddr(host1)
	existingAddrIPVer := addr1.Is4()

	addr2, _ := netip.ParseAddr(host2)
	newAddrIPVer := addr2.Is4()

	return existingAddrIPVer == newAddrIPVer
}

// validateAdvertisedListeners performs common checks on the advertised listers as per Kafka specification https://kafka.apache.org/documentation/#brokerconfigs_advertised.listeners.
// Unlike with listeners, having duplicated ports is allowed. The only constraint is advertising to 0.0.0.0 is not allowed.
func validateAdvertisedListeners(b *config.Broker) error {
	for _, listener := range b.AdvertisedListeners {
		if strings.EqualFold(listener.Host, "0.0.0.0") || listener.Host == "" {
			return fmt.Errorf("advertising listener on 0.0.0.0 address is not allowed for listener %s", listener.ListenerName)
		}
	}

	return nil
}
