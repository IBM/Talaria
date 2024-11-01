// protocol has been generated from message format json - DO NOT EDIT
package protocol

// PartitionData_BeginQuorumEpochRequest contains a
type PartitionData_BeginQuorumEpochRequest struct {
	// Version defines the protocol version to use for encode and decode
	Version int16
	// PartitionIndex contains the partition index.
	PartitionIndex int32
	// LeaderID contains the ID of the newly elected leader
	LeaderID int32
	// LeaderEpoch contains the epoch of the newly elected leader
	LeaderEpoch int32
}

func (p *PartitionData_BeginQuorumEpochRequest) encode(pe packetEncoder, version int16) (err error) {
	p.Version = version
	pe.putInt32(p.PartitionIndex)

	pe.putInt32(p.LeaderID)

	pe.putInt32(p.LeaderEpoch)

	return nil
}

func (p *PartitionData_BeginQuorumEpochRequest) decode(pd packetDecoder, version int16) (err error) {
	p.Version = version
	if p.PartitionIndex, err = pd.getInt32(); err != nil {
		return err
	}

	if p.LeaderID, err = pd.getInt32(); err != nil {
		return err
	}

	if p.LeaderEpoch, err = pd.getInt32(); err != nil {
		return err
	}

	return nil
}

// TopicData_BeginQuorumEpochRequest contains a
type TopicData_BeginQuorumEpochRequest struct {
	// Version defines the protocol version to use for encode and decode
	Version int16
	// TopicName contains the topic name.
	TopicName string
	// Partitions contains a
	Partitions []PartitionData_BeginQuorumEpochRequest
}

func (t *TopicData_BeginQuorumEpochRequest) encode(pe packetEncoder, version int16) (err error) {
	t.Version = version
	if err := pe.putString(t.TopicName); err != nil {
		return err
	}

	if err := pe.putArrayLength(len(t.Partitions)); err != nil {
		return err
	}
	for _, block := range t.Partitions {
		if err := block.encode(pe, t.Version); err != nil {
			return err
		}
	}

	return nil
}

func (t *TopicData_BeginQuorumEpochRequest) decode(pd packetDecoder, version int16) (err error) {
	t.Version = version
	if t.TopicName, err = pd.getString(); err != nil {
		return err
	}

	var numPartitions int
	if numPartitions, err = pd.getArrayLength(); err != nil {
		return err
	}
	if numPartitions > 0 {
		t.Partitions = make([]PartitionData_BeginQuorumEpochRequest, numPartitions)
		for i := 0; i < numPartitions; i++ {
			var block PartitionData_BeginQuorumEpochRequest
			if err := block.decode(pd, t.Version); err != nil {
				return err
			}
			t.Partitions[i] = block
		}
	}

	return nil
}

type BeginQuorumEpochRequest struct {
	// Version defines the protocol version to use for encode and decode
	Version int16
	// ClusterID contains a
	ClusterID *string
	// Topics contains a
	Topics []TopicData_BeginQuorumEpochRequest
}

func (r *BeginQuorumEpochRequest) encode(pe packetEncoder) (err error) {
	if err := pe.putNullableString(r.ClusterID); err != nil {
		return err
	}

	if err := pe.putArrayLength(len(r.Topics)); err != nil {
		return err
	}
	for _, block := range r.Topics {
		if err := block.encode(pe, r.Version); err != nil {
			return err
		}
	}

	return nil
}

func (r *BeginQuorumEpochRequest) decode(pd packetDecoder, version int16) (err error) {
	r.Version = version
	if r.ClusterID, err = pd.getNullableString(); err != nil {
		return err
	}

	var numTopics int
	if numTopics, err = pd.getArrayLength(); err != nil {
		return err
	}
	if numTopics > 0 {
		r.Topics = make([]TopicData_BeginQuorumEpochRequest, numTopics)
		for i := 0; i < numTopics; i++ {
			var block TopicData_BeginQuorumEpochRequest
			if err := block.decode(pd, r.Version); err != nil {
				return err
			}
			r.Topics[i] = block
		}
	}

	return nil
}

func (r *BeginQuorumEpochRequest) GetKey() int16 {
	return 53
}

func (r *BeginQuorumEpochRequest) GetVersion() int16 {
	return r.Version
}

func (r *BeginQuorumEpochRequest) GetHeaderVersion() int16 {
	return 1
}

func (r *BeginQuorumEpochRequest) IsValidVersion() bool {
	return r.Version == 0
}

func (r *BeginQuorumEpochRequest) GetRequiredVersion() int16 {
	// TODO - it isn't possible to determine this from the message format json files
	return 0
}
