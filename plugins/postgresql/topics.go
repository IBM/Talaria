package postgresql

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/ibm/opentalaria/protocol"
)

func (p *Plugin) AddTopic(topic protocol.CreatableTopic) error {
	slog.Debug("Create", "topic", topic.Name, "configs", topic.Configs)

	statement := `
		INSERT INTO topics (topic_id, topic_name, num_partitions, replication_factor)
		VALUES ($1, $2, $3, $4)`

	_, err := p.db.Exec(statement, uuid.New(), topic.Name, topic.NumPartitions, topic.ReplicationFactor)

	return err
}
