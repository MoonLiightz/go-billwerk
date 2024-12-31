package optimize

import (
	"context"
	"fmt"
	"time"
)

// PlanState represents the state of a subscription plan.
type PlanState string

const (
	PlanStateActive     PlanState = "active"     // Plan is currently active.
	PlanStateSuperseded PlanState = "superseded" // Plan has been replaced by another version.
	PlanStateDeleted    PlanState = "deleted"    // Plan is deleted.
)

// PlanScheduleType defines the scheduling type for a plan.
type PlanScheduleType string

const (
	PlanScheduleTypeManual         PlanScheduleType = "manual"          // Manually scheduled plans.
	PlanScheduleTypeDaily          PlanScheduleType = "daily"           // Daily scheduled plans.
	PlanScheduleTypeWeeklyFixedDay PlanScheduleType = "weekly_fixedday" // Weekly scheduling on fixed days.
	PlanScheduleTypeMonthStartDate PlanScheduleType = "month_startdate" // Monthly scheduling based on start date.
)

// PlanPartialPeriodHandling defines how to handle partial billing periods.
type PlanPartialPeriodHandling string

const (
	PlanPartialPeriodHandlingBillFull       PlanPartialPeriodHandling = "bill_full"        // Bill full amount.
	PlanPartialPeriodHandlingBillProrated   PlanPartialPeriodHandling = "bill_prorated"    // Bill prorated amount.
	PlanPartialPeriodHandlingBillZeroAmount PlanPartialPeriodHandling = "bill_zero_amount" // Bill zero amount.
	PlanPartialPeriodHandlingNoBill         PlanPartialPeriodHandling = "no_bill"          // No billing.
)

// PlanSetupFeeHandling specifies how to handle setup fees.
type PlanSetupFeeHandling string

const (
	First               PlanSetupFeeHandling = "first"                // Include setup fee on the first invoice.
	Separate            PlanSetupFeeHandling = "separate"             // Separate invoice for the setup fee.
	SeparateConditional PlanSetupFeeHandling = "separate_conditional" // Conditional separate setup fee billing.
)

// PlanFixedLifeTimeUnit represents time units for fixed life plans.
type PlanFixedLifeTimeUnit string

const (
	PlanFixedLifeTimeUnitDays   PlanFixedLifeTimeUnit = "days"   // Life time in days.
	PlanFixedLifeTimeUnitMonths PlanFixedLifeTimeUnit = "months" // Life time in months.
)

// PlanTrialIntervalUnit defines the time unit for a trial period.
type PlanTrialIntervalUnit string

const (
	PlanTrialIntervalUnitDays   PlanTrialIntervalUnit = "days"   // Trial period in days.
	PlanTrialIntervalUnitMonths PlanTrialIntervalUnit = "months" // Trial period in months.
)

// PlanSupersedeMode represents the mode of superseding a plan.
type PlanSupersedeMode string

const (
	NoSubUpdate        PlanSupersedeMode = "no_sub_update"        // No subscription update.
	ScheduledSubUpdate PlanSupersedeMode = "scheduled_sub_update" // Subscription update scheduled.
)

// PlanRange represents the range for retrieving plans.
type PlanRange string

const (
	PlanRangeCreated PlanRange = "created" // Retrieve plans by creation date.
)

// Plan defines the structure for a subscription plan.
type Plan struct {
	// Name of the plan.
	Name string `json:"name"`

	// Description of the plan.
	Description string `json:"description,omitempty"`

	// Optional vat for this plan. Account default is used if none given.
	Vat float64 `json:"vat,omitempty"`

	// Amount for the plan in the smallest unit for the account currency.
	Amount int32 `json:"amount"`

	// Optional default quantity of the subscription plan product for new subscriptions. Default is 1.
	Quantity int32 `json:"quantity,omitempty"`

	// Subscriptions can either be prepaid where an amount is paid in advance, or the opposite.
	// This setting only relates to handling of pause scenarios.
	Prepaid bool `json:"prepaid,omitempty"`

	// Per account unique handle for the subscription plan. Max length 255 with allowable characters [a-zA-Z0-9_.-@].
	Handle string `json:"handle"`

	// Plan version
	Version int32 `json:"version,omitempty"`

	// State of the subscription plan one of the following: active, superseded, deleted
	State PlanState `json:"state,omitempty"`

	// Currency for the subscription plan in ISO 4217 three letter alpha code.
	Currency string `json:"currency,omitempty"`

	// Date when the subscription plan was created. In ISO-8601 extended offset date-time format.
	Created *time.Time `json:"created,omitempty"`

	// Date when the subscription plan was deleted. In ISO-8601 extended offset date-time format.
	Deleted *time.Time `json:"deleted,omitempty"`

	// Dunning plan by handle to use for the subscription plan. Default dunning plan will be used if none given.
	DunningPlan string `json:"dunning_plan,omitempty"`

	// Optional tax policy handle for this plan. If vat and tax policy is given, vat will be ignored.
	TaxPolicy string `json:"tax_policy,omitempty"`

	// Optional renewal reminder email settings. Number of days before next billing to send a reminder email.
	RenewalReminderEmailDays int32 `json:"renewal_reminder_email_days,omitempty"`

	// Optional end of trial reminder email settings. Number of days before end of trial to send a reminder email.
	TrialReminderEmailDays int32 `json:"trial_reminder_email_days,omitempty"`

	// How to handle a potential initial partial billing period for fixed day scheduling.
	// The options are to bill for a full period, bill prorated for the partial period, bill a zero amoumt,
	// or not to consider the period before first fixed day a billing period.
	// The default is to bill prorated. Options: bill_full, bill_prorated, bill_zero_amount, no_bill.
	PartialPeriodHandling PlanPartialPeriodHandling `json:"partial_period_handling,omitempty"`

	// Whether to add a zero amount order line to subscription invoices if plan amount is zero
	// or the subscription overrides to zero amount. The default is to not include the line.
	// If no other order lines are present the plan order line will be added.
	IncludeZeroAmount bool `json:"include_zero_amount,omitempty"`

	// Optional one-time setup fee billed with the first invoice or as a separate invoice
	// depending on the setting setup_fee_handling.
	SetupFee int32 `json:"setup_fee,omitempty"`

	// Optional invoice order text for the setup fee that.
	SetupFeeText string `json:"setup_fee_text,omitempty"`

	// How the billing of the setup fee should be done.
	// The options are:
	// first - include setup fee as order line on the first scheduled invoice.
	// separate - create a separate invoice for the setup fee, is appropriate if first invoice is not in conjunction with creation.
	// separate_conditional - create a separate invoice for setup fee if the first invoice is not created in conjunction with the creation.
	// Default is first.
	SetupFeeHandling string `json:"setup_fee_handling,omitempty"`

	// For fixed day scheduling and prorated partial handling calculate prorated amount
	// using whole days counting start day as a full day, or use by the minute proration
	// calculation from start date time to the next period start.
	// Default is true (whole days).
	PartialProrationDays bool `json:"partial_proration_days,omitempty"`

	// When using trial for fixed day scheduling use this setting to control if trial
	// expires at midnight or the trial period is down to the minute.
	// Default is true (trial until start of day).
	// Trial in days can only be true if partial_proration_days is also set to true.
	FixedTrialDays bool `json:"fixed_trial_days,omitempty"`

	// When using prorated partial handling the prorated amount for plan and add-ons might result in very small amounts.
	// A minimum prorated amount for plan and add-ons can be defined.
	// If the prorated amount is below this minimum the amount will be changed to zero.
	MinimumProratedAmount int32 `json:"minimum_prorated_amount,omitempty"`

	// Indicates that Account Funding Transaction (AFT) is requested.
	AccountFunding bool `json:"account_funding,omitempty"`

	// Whether the amount is including VAT. Default true.
	AmountInclVat bool `json:"amount_incl_vat,omitempty"`

	// Fixed number of renewals for subscriptions using this plan. Equals the number of scheduled invoices.
	// Default is no fixed amount of renewals.
	FixedCount int32 `json:"fixed_count,omitempty"`

	// Time unit use for fixed life time (months, days).
	FixedLifeTimeUnit string `json:"fixed_life_time_unit,omitempty"`

	// Optional fixed life time length for subscriptions using this plan. E.g. 12 months.
	// Subscriptions will cancel after the fixed life time and expire when the active billing cycle ends.
	FixedLifeTimeLength int32 `json:"fixed_life_time_length,omitempty"`

	// Time unit for free trial period (months, days).
	TrialIntervalUnit PlanTrialIntervalUnit `json:"trial_interval_unit,omitempty"`

	// Optional free trial interval length. E.g. 1 months.
	TrialIntervalLength int32 `json:"trial_interval_length,omitempty"`

	// The length of intervals. E.g. every second month or every 14 days.
	IntervalLength int32 `json:"interval_length,omitempty"`

	// Scheduling type, one of the following:
	// manual, daily, weekly_fixedday, month_startdate, month_fixedday, month_lastday.
	// See documentation for descriptions of the different types.
	ScheduleType PlanScheduleType `json:"schedule_type"`

	// If a fixed day scheduling type is used a fixed day must be provided.
	// For months the allowed value is 1-28 for weeks it is 1-7.
	ScheduleFixedDay int32 `json:"schedule_fixed_day,omitempty"`

	// For fixed month schedule types the base month can be used to control
	// which months are eligible for start of first billing period.
	// The eligible months are calculated as base_month + k * interval_length up to 12.
	// E.g. to use quaterly billing in the months jan-apr-jul-oct, base_month 1 and interval_length 3 can be used.
	// If not defined the first fixed day will be used as start of first billing period.
	BaseMonth int32 `json:"base_month,omitempty"`

	// Optional number of notice periods for a cancel.
	// The subscription will be cancelled for this number of full periods before expiring.
	// Either from the cancellation date, or from the end of the current period.
	// See notice_periods_after_current. The default is to expire at the end of current period (0).
	// A value of 1 (and notice_periods_after_current set to true) will for example result
	// in a scenario where the subscription is cancelled until the end of current period,
	// and then for the full subsequent period before expiring.
	NoticePeriods int32 `json:"notice_periods,omitempty"`

	// If notice periods is set, this option controls whether the number of full notice periods
	// should start at the end of the current period, or run from cancellation date and result in a partial period
	// with partial amount for the last period.
	// The default is true.
	// E.g. if set to false and notice_periods = 1 then the subscription will be cancelled for exactly for one period
	// from the cancellation time and a partial amount will be billed at the start of the next billing period.
	NoticePeriodsAfterCurrent bool `json:"notice_periods_after_current,omitempty"`

	// Optional number of fixation periods. Fixation periods will guarantee that a subscription
	// will have this number of paid full periods before expiring after a cancel.
	// Default is to have no requirement (0).
	FixationPeriods int32 `json:"fixation_periods,omitempty"`

	// If fixation periods are defined, and the subscription can have a partial prorated first period,
	// this parameter controls if the last period should be full,
	// or partial to give exactly fixation_periods paid periods.
	// Default is false.
	FixationPeriodsFull bool `json:"fixation_periods_full,omitempty"`

	// List of entitlement handles to be added to the plan.
	Entitlements []string `json:"entitlements,omitempty"`
}

// PlanSupersede includes additional fields for superseding a plan.
type PlanSupersede struct {
	Plan
	SupersedeMode PlanSupersedeMode `json:"supersede_mode,omitempty"` // Supersede mode for the plan.
}

// ListOfPlansResponse contains the response for listing plans.
type ListOfPlansResponse struct {
	Size          int       `json:"size"`            // Number of plans returned.
	Count         int       `json:"count"`           // Total count of plans.
	To            string    `json:"to"`              // End of the range.
	From          string    `json:"from"`            // Start of the range.
	Content       []*Plan   `json:"content"`         // List of plans.
	Range         PlanRange `json:"range"`           // Plan range.
	NextPageToken string    `json:"next_page_token"` // Token for the next page of results.
}

// PlanEntitlement defines entitlements associated with a plan.
type PlanEntitlement struct {
	Handle      string     `json:"handle"`      // Unique handle for the entitlement.
	Name        string     `json:"name"`        // Name of the entitlement.
	Description string     `json:"description"` // Description of the entitlement.
	Created     *time.Time `json:"created"`     // Creation date of the entitlement.
}

// GetListOfPlans retrieves a list of plans based on the provided query parameters.
func (b *Billwerk) GetListOfPlans(ctx context.Context, params ...QueryParamFunc) (*ListOfPlansResponse, error) {
	endpoint := "/list/plan"

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint)

	for _, param := range params {
		param(requestBuilder)
	}

	req, err := requestBuilder.GET()
	if err != nil {
		return nil, err
	}

	var res ListOfPlansResponse
	if err = b.Do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetPlan retrieves a specific plan by its handle.
func (b *Billwerk) GetPlan(ctx context.Context, handle string, params ...QueryParamFunc) (*Plan, error) {
	endpoint := fmt.Sprintf("/plan/%s/current", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint)

	for _, param := range params {
		param(requestBuilder)
	}

	req, err := requestBuilder.GET()
	if err != nil {
		return nil, err
	}

	var res Plan
	if err = b.Do(req, &res); err != nil {
		fmt.Printf("%+v\n", err)
		return nil, err
	}

	return &res, nil
}

// GetListOfPlanVersions retrieves all versions of a plan by its handle.
func (b *Billwerk) GetListOfPlanVersions(ctx context.Context, handle string, params ...QueryParamFunc) ([]*Plan, error) {
	endpoint := fmt.Sprintf("/plan/%s", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint)

	for _, param := range params {
		param(requestBuilder)
	}

	req, err := requestBuilder.GET()
	if err != nil {
		return nil, err
	}

	var res []*Plan
	if err = b.Do(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// CreatePlan creates a new subscription plan.
func (b *Billwerk) CreatePlan(ctx context.Context, plan *Plan) (*Plan, error) {
	endpoint := "/plan"

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint).
		WithJSONBody(plan)

	req, err := requestBuilder.POST()
	if err != nil {
		return nil, err
	}

	var res Plan
	if err = b.Do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// SupersedePlan supersedes an existing plan with a new version.
func (b *Billwerk) SupersedePlan(ctx context.Context, handle string, plan *PlanSupersede) (*Plan, error) {
	endpoint := fmt.Sprintf("/plan/%s", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint).
		WithJSONBody(plan)

	req, err := requestBuilder.POST()
	if err != nil {
		return nil, err
	}

	var res Plan
	if err = b.Do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UpdatePlan updates an existing subscription plan by its handle.
func (b *Billwerk) UpdatePlan(ctx context.Context, handle string, plan *Plan) (*Plan, error) {
	endpoint := fmt.Sprintf("/plan/%s", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint).
		WithJSONBody(plan)

	req, err := requestBuilder.PUT()
	if err != nil {
		return nil, err
	}

	var res Plan
	if err = b.Do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// DeletePlan deletes a subscription plan by its handle.
func (b *Billwerk) DeletePlan(ctx context.Context, handle string) (*Plan, error) {
	endpoint := fmt.Sprintf("/plan/%s", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint)

	req, err := requestBuilder.DELETE()
	if err != nil {
		return nil, err
	}

	var res Plan
	if err = b.Do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// UndeletePlan undeletes a previously deleted subscription plan by its handle.
func (b *Billwerk) UndeletePlan(ctx context.Context, handle string) (*Plan, error) {
	endpoint := fmt.Sprintf("/plan/%s/undelete", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint)

	req, err := requestBuilder.POST()
	if err != nil {
		return nil, err
	}

	var res Plan
	if err = b.Do(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetPlanEntitlements retrieves entitlements associated with a specific plan version by its handle and version.
func (b *Billwerk) GetPlanEntitlements(ctx context.Context, handle string, version int32) ([]*PlanEntitlement, error) {
	endpoint := fmt.Sprintf("/plan/%s/%d/entitlement", handle, version)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint)

	req, err := requestBuilder.GET()
	if err != nil {
		return nil, err
	}

	var res []*PlanEntitlement
	if err = b.Do(req, &res); err != nil {
		return nil, err
	}

	return res, nil
}

// GetPlanMetadata retrieves the metadata for a plan by its handle.
// The result is stored in the metadata parameter and should be a pointer e.g. &map[string]interface{}{}
// or &struct{}{} with the expected fields / json tags.
func (b *Billwerk) GetPlanMetadata(ctx context.Context, handle string, metadata interface{}) error {
	endpoint := fmt.Sprintf("/plan/%s/metadata", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint)

	req, err := requestBuilder.GET()
	if err != nil {
		return err
	}

	if err = b.Do(req, &metadata); err != nil {
		return err
	}

	return nil
}

// CreateOrUpdatePlanMetadata creates or updates the metadata for a plan by its handle.
// The response is stored in the metadata parameter and modifies the passed in object.
func (b *Billwerk) CreateOrUpdatePlanMetadata(ctx context.Context, handle string, metadata interface{}) error {
	endpoint := fmt.Sprintf("/plan/%s/metadata", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint).
		WithJSONBody(metadata)

	req, err := requestBuilder.PUT()
	if err != nil {
		return err
	}

	if err = b.Do(req, metadata); err != nil {
		return err
	}

	return nil
}

// DeletePlanMetadata deletes metadata associated with a specific plan by its handle.
func (b *Billwerk) DeletePlanMetadata(ctx context.Context, handle string) error {
	endpoint := fmt.Sprintf("/plan/%s/metadata", handle)

	requestBuilder := b.newBillwerkRequest(ctx).
		WithEndpoint(endpoint)

	req, err := requestBuilder.DELETE()
	if err != nil {
		return err
	}

	if err = b.Do(req, nil); err != nil {
		return err
	}

	return nil
}
