package main

import (
	"time"

	"github.com/guregu/null"
)

// EffectOutput is a representation of an operation that aligns with the BigQuery table history_effects
type EffectOutput struct {
	Address        string                 `json:"address"`
	AddressMuxed   null.String            `json:"address_muxed,omitempty"`
	OperationID    int64                  `json:"operation_id"`
	Details        map[string]interface{} `json:"details"`
	Type           int32                  `json:"type"`
	TypeString     string                 `json:"type_string"`
	LedgerClosed   time.Time              `json:"closed_at"`
	LedgerSequence uint32                 `json:"ledger_sequence"`
	EffectIndex    uint32                 `json:"index"`
	EffectId       string                 `json:"id"`
}

// EffectType is the numeric type for an effect
type EffectType int

// Effect type constants - these are directly from the Stellar effects code
const (
	EffectAccountCreated                     EffectType = 0
	EffectAccountRemoved                     EffectType = 1
	EffectAccountCredited                    EffectType = 2
	EffectAccountDebited                     EffectType = 3
	EffectAccountThresholdsUpdated           EffectType = 4
	EffectAccountHomeDomainUpdated           EffectType = 5
	EffectAccountFlagsUpdated                EffectType = 6
	EffectAccountInflationDestinationUpdated EffectType = 7
	EffectSignerCreated                      EffectType = 10
	EffectSignerRemoved                      EffectType = 11
	EffectSignerUpdated                      EffectType = 12
	EffectTrustlineCreated                   EffectType = 20
	EffectTrustlineRemoved                   EffectType = 21
	EffectTrustlineUpdated                   EffectType = 22
	EffectTrustlineFlagsUpdated              EffectType = 26
	EffectOfferCreated                       EffectType = 30
	EffectOfferRemoved                       EffectType = 31
	EffectOfferUpdated                       EffectType = 32
	EffectTrade                              EffectType = 33
	EffectDataCreated                        EffectType = 40
	EffectDataRemoved                        EffectType = 41
	EffectDataUpdated                        EffectType = 42
	EffectSequenceBumped                     EffectType = 43
	EffectClaimableBalanceCreated            EffectType = 50
	EffectClaimableBalanceClaimantCreated    EffectType = 51
	EffectClaimableBalanceClaimed            EffectType = 52
	EffectAccountSponsorshipCreated          EffectType = 60
	EffectAccountSponsorshipUpdated          EffectType = 61
	EffectAccountSponsorshipRemoved          EffectType = 62
	EffectTrustlineSponsorshipCreated        EffectType = 63
	EffectTrustlineSponsorshipUpdated        EffectType = 64
	EffectTrustlineSponsorshipRemoved        EffectType = 65
	EffectDataSponsorshipCreated             EffectType = 66
	EffectDataSponsorshipUpdated             EffectType = 67
	EffectDataSponsorshipRemoved             EffectType = 68
	EffectClaimableBalanceSponsorshipCreated EffectType = 69
	EffectClaimableBalanceSponsorshipUpdated EffectType = 70
	EffectClaimableBalanceSponsorshipRemoved EffectType = 71
	EffectSignerSponsorshipCreated           EffectType = 72
	EffectSignerSponsorshipUpdated           EffectType = 73
	EffectSignerSponsorshipRemoved           EffectType = 74
	EffectClaimableBalanceClawedBack         EffectType = 80
	EffectLiquidityPoolDeposited             EffectType = 90
	EffectLiquidityPoolWithdrew              EffectType = 91
	EffectLiquidityPoolTrade                 EffectType = 92
	EffectLiquidityPoolCreated               EffectType = 93
	EffectLiquidityPoolRemoved               EffectType = 94
	EffectLiquidityPoolRevoked               EffectType = 95
	EffectContractCredited                   EffectType = 96
	EffectContractDebited                    EffectType = 97
	EffectExtendFootprintTtl                 EffectType = 98
	EffectRestoreFootprint                   EffectType = 99
)

// EffectTypeNames stores a map of effect type ID and names
var EffectTypeNames = map[EffectType]string{
	EffectAccountCreated:                     "account_created",
	EffectAccountRemoved:                     "account_removed",
	EffectAccountCredited:                    "account_credited",
	EffectAccountDebited:                     "account_debited",
	EffectAccountThresholdsUpdated:           "account_thresholds_updated",
	EffectAccountHomeDomainUpdated:           "account_home_domain_updated",
	EffectAccountFlagsUpdated:                "account_flags_updated",
	EffectAccountInflationDestinationUpdated: "account_inflation_destination_updated",
	EffectSignerCreated:                      "signer_created",
	EffectSignerRemoved:                      "signer_removed",
	EffectSignerUpdated:                      "signer_updated",
	EffectTrustlineCreated:                   "trustline_created",
	EffectTrustlineRemoved:                   "trustline_removed",
	EffectTrustlineUpdated:                   "trustline_updated",
	EffectTrustlineFlagsUpdated:              "trustline_flags_updated",
	EffectOfferCreated:                       "offer_created",
	EffectOfferRemoved:                       "offer_removed",
	EffectOfferUpdated:                       "offer_updated",
	EffectTrade:                              "trade",
	EffectDataCreated:                        "data_created",
	EffectDataRemoved:                        "data_removed",
	EffectDataUpdated:                        "data_updated",
	EffectSequenceBumped:                     "sequence_bumped",
	EffectClaimableBalanceCreated:            "claimable_balance_created",
	EffectClaimableBalanceClaimed:            "claimable_balance_claimed",
	EffectClaimableBalanceClaimantCreated:    "claimable_balance_claimant_created",
	EffectAccountSponsorshipCreated:          "account_sponsorship_created",
	EffectAccountSponsorshipUpdated:          "account_sponsorship_updated",
	EffectAccountSponsorshipRemoved:          "account_sponsorship_removed",
	EffectTrustlineSponsorshipCreated:        "trustline_sponsorship_created",
	EffectTrustlineSponsorshipUpdated:        "trustline_sponsorship_updated",
	EffectTrustlineSponsorshipRemoved:        "trustline_sponsorship_removed",
	EffectDataSponsorshipCreated:             "data_sponsorship_created",
	EffectDataSponsorshipUpdated:             "data_sponsorship_updated",
	EffectDataSponsorshipRemoved:             "data_sponsorship_removed",
	EffectClaimableBalanceSponsorshipCreated: "claimable_balance_sponsorship_created",
	EffectClaimableBalanceSponsorshipUpdated: "claimable_balance_sponsorship_updated",
	EffectClaimableBalanceSponsorshipRemoved: "claimable_balance_sponsorship_removed",
	EffectSignerSponsorshipCreated:           "signer_sponsorship_created",
	EffectSignerSponsorshipUpdated:           "signer_sponsorship_updated",
	EffectSignerSponsorshipRemoved:           "signer_sponsorship_removed",
	EffectClaimableBalanceClawedBack:         "claimable_balance_clawed_back",
	EffectLiquidityPoolDeposited:             "liquidity_pool_deposited",
	EffectLiquidityPoolWithdrew:              "liquidity_pool_withdrew",
	EffectLiquidityPoolTrade:                 "liquidity_pool_trade",
	EffectLiquidityPoolCreated:               "liquidity_pool_created",
	EffectLiquidityPoolRemoved:               "liquidity_pool_removed",
	EffectLiquidityPoolRevoked:               "liquidity_pool_revoked",
	EffectContractCredited:                   "contract_credited",
	EffectContractDebited:                    "contract_debited",
	EffectExtendFootprintTtl:                 "extend_footprint_ttl",
	EffectRestoreFootprint:                   "restore_footprint",
}
