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

// Initialize processes configuration parameters.
func (p *EffectsProcessor) Initialize(config map[string]interface{}) error {
	p.config = config

	// Extract network passphrase from config
	if passphrase, ok := config["network_passphrase"].(string); ok {
		p.networkPassphrase = passphrase
	} else {
		return errors.New("network_passphrase is required in configuration")
	}

	log.Println("EffectsProcessor initialized with config:", config)
	return nil
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
	if err := json.Unmarshal(msg.Payload(), &transaction); err != nil {
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
		outputMsg := pluginapi.NewMessage(effectJSON)

		// Set message metadata from the original message
		for k, v := range msg.Metadata() {
			outputMsg.SetMetadata(k, v)
		}

		// Add effect-specific metadata
		outputMsg.SetMetadata("effect_id", effect.EffectId)
		outputMsg.SetMetadata("effect_type", effect.TypeString)

		// Output the message
		if err := msg.OutputTo(ctx, outputMsg); err != nil {
			return NewProcessorError(
				fmt.Errorf("error outputting effect message: %w", err),
				ErrorTypeIO,
				ErrorSeverityError,
			)
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
