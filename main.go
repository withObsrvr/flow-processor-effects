package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/withObsrvr/pluginapi"
)

// EffectsProcessor is a Flow processor plugin that transforms operations into effects.
type EffectsProcessor struct {
	config            map[string]interface{}
	networkPassphrase string
	consumers         []pluginapi.Consumer
}

// Name returns the plugin's name.
func (p *EffectsProcessor) Name() string {
	return "flow/processor/effects"
}

// Version returns the plugin version.
func (p *EffectsProcessor) Version() string {
	return "0.1.0"
}

// Type returns the type of the plugin.
func (p *EffectsProcessor) Type() pluginapi.PluginType {
	return pluginapi.ProcessorPlugin
}

// GetSchemaDefinition returns GraphQL type definitions for this plugin
func (p *EffectsProcessor) GetSchemaDefinition() string {
	return `
type Effect {
    id: String!
    address: String!
    addressMuxed: String
    operationId: Int!
    details: JSON
    type: Int!
    typeString: String!
    closedAt: String!
    ledgerSequence: Int!
    index: Int!
}

scalar JSON
`
}

// GetQueryDefinitions returns GraphQL query definitions for this plugin
func (p *EffectsProcessor) GetQueryDefinitions() string {
	return `
    effectsByOperationId(operationId: Int!): [Effect]
    effectsByAddress(address: String!): [Effect]
    effectsByType(type: Int!): [Effect]
`
}

// Initialize processes configuration parameters.
func (p *EffectsProcessor) Initialize(config map[string]interface{}) error {
	p.config = config
	p.consumers = make([]pluginapi.Consumer, 0)

	// Extract network passphrase from config
	if passphrase, ok := config["network_passphrase"].(string); ok {
		p.networkPassphrase = passphrase
	} else {
		return errors.New("network_passphrase is required in configuration")
	}

	log.Println("EffectsProcessor initialized with config:", config)
	return nil
}

// RegisterConsumer registers a downstream consumer
func (p *EffectsProcessor) RegisterConsumer(consumer pluginapi.Consumer) {
	log.Printf("EffectsProcessor: Registering consumer %s", consumer.Name())
	p.consumers = append(p.consumers, consumer)
}

// Process implements the main processing logic for transforming operations into effects.
func (p *EffectsProcessor) Process(ctx context.Context, msg pluginapi.Message) error {
	// Check for context cancellation
	if err := ctx.Err(); err != nil {
		return NewProcessorError(
			fmt.Errorf("context canceled during processing: %w", err),
			ErrorTypeProcessing,
			ErrorSeverityWarning,
		)
	}

	// Extract transaction data from the message
	var transaction map[string]interface{}

	// Use type assertion to convert payload to []byte
	payloadBytes, ok := msg.Payload.([]byte)
	if !ok {
		return NewProcessorError(
			fmt.Errorf("expected payload to be []byte, got %T", msg.Payload),
			ErrorTypeParsing,
			ErrorSeverityError,
		)
	}

	if err := json.Unmarshal(payloadBytes, &transaction); err != nil {
		return NewProcessorError(
			fmt.Errorf("error unmarshaling transaction: %w", err),
			ErrorTypeParsing,
			ErrorSeverityError,
		)
	}

	// Process the transaction and generate effects
	effects, err := p.transformTransactionToEffects(ctx, transaction)
	if err != nil {
		return NewProcessorError(
			fmt.Errorf("error transforming transaction to effects: %w", err),
			ErrorTypeProcessing,
			ErrorSeverityError,
		).WithTransaction(getTransactionHash(transaction))
	}

	// If no effects, just return
	if len(effects) == 0 {
		return nil
	}

	// Convert effects to messages and output them
	for _, effect := range effects {
		effectJSON, err := json.Marshal(effect)
		if err != nil {
			return NewProcessorError(
				fmt.Errorf("error marshaling effect: %w", err),
				ErrorTypeParsing,
				ErrorSeverityWarning,
			)
		}

		// Create a new message with the effect data
		outputMsg := pluginapi.Message{
			Payload: effectJSON,
			Metadata: map[string]interface{}{
				"effect_id":   effect.EffectId,
				"effect_type": effect.TypeString,
			},
		}

		// Copy existing metadata from the source message
		for k, v := range msg.Metadata {
			outputMsg.Metadata[k] = v
		}

		// Forward to consumers
		for _, consumer := range p.consumers {
			if err := consumer.Process(ctx, outputMsg); err != nil {
				log.Printf("Error in consumer %s: %v", consumer.Name(), err)
			}
		}
	}

	return nil
}

// Close handles any cleanup if necessary.
func (p *EffectsProcessor) Close() error {
	log.Println("EffectsProcessor closed")
	return nil
}

// Helper function to get transaction hash
func getTransactionHash(transaction map[string]interface{}) string {
	if hash, ok := transaction["hash"].(string); ok {
		return hash
	}
	return "unknown"
}

// New is the exported function for dynamic plugin loading.
func New() pluginapi.Plugin {
	return &EffectsProcessor{}
}
