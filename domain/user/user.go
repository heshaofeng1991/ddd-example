package user

import (
	"time"

	domainTenant "github.com/heshaofeng1991/ddd-johnny/domain/tenant"
	"github.com/pkg/errors"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrEmailExists            = errors.New("email already exists")
	ErrInvalidEmailOrPassword = errors.New("Invalid email or password")
)

type User struct {
	userID   int64
	username string
	email    string
	password string
	phone    string

	website  string
	platform string
	concerns string

	source    string
	sourceTag string
	status    int8

	storeCode string

	guideInfo *GuideInfo

	lastLoggedTime time.Time
	createdAt      time.Time
	updatedAt      time.Time

	tenant *domainTenant.Tenant
}

func (u *User) Website() string {
	return u.website
}

func (u *User) Platform() string {
	return u.platform
}

func (u *User) Concerns() string {
	return u.concerns
}

func (u *User) GuideInfo() *GuideInfo {
	return u.guideInfo
}

func (u *User) StoreCode() string {
	return u.storeCode
}

const (
	GuideStep1 = "signup"
	GuideStep2 = "info"
	GuideStep3 = "integration"
)

type GuideInfo struct {
	status     int
	finished   bool
	HsObjectID string
	steps      []*Steps
	questions  []*Question
}

func (g *GuideInfo) Steps() []*Steps {
	return g.steps
}

func (g *GuideInfo) SetSteps(steps []*Steps) {
	g.steps = steps
}

func (g *GuideInfo) SetStatus(status int) {
	g.status = status
}

func (g *GuideInfo) SetFinished(finished bool) {
	g.finished = finished
}

func (g *GuideInfo) SetQuestions(questions []*Question) {
	g.questions = questions
}

func (g GuideInfo) Status() int {
	return g.status
}

func (g GuideInfo) Finished() bool {
	return g.finished
}

func (g GuideInfo) Questions() []*Question {
	return g.questions
}

func NewGuideInfo(status int, finished bool, hsObjectID string, questions []*Question) *GuideInfo {
	return &GuideInfo{
		status:     status,
		finished:   finished,
		HsObjectID: hsObjectID,
		questions:  questions,
		steps: []*Steps{
			{Title: GuideStep1, Step: 1},
			{Title: GuideStep2, Step: 2}, //nolint:gomnd
			{Title: GuideStep3, Step: 3}, //nolint:gomnd
		},
	}
}

type Steps struct {
	Title  string
	Step   int
	Status string
}

type Question struct {
	Title  string
	Answer interface{}
}

// Option is a function that can be used as an option when creating a new user.
type Option func(u *User)

func WithGuideInfo(info *GuideInfo) Option {
	return func(u *User) {
		u.guideInfo = info
	}
}

func WithWebsite(website *string) Option {
	return func(u *User) {
		if website != nil {
			u.website = *website
		}
	}
}

func WithStoreCode(storeCode string) Option {
	return func(u *User) {
		u.storeCode = storeCode
	}
}

func WithPlatform(platform *string) Option {
	return func(u *User) {
		if platform != nil {
			u.platform = *platform
		}
	}
}

func WithConcerns(concerns *string) Option {
	return func(u *User) {
		if concerns != nil {
			u.concerns = *concerns
		}
	}
}

func WithSource(source *string) Option {
	return func(u *User) {
		if source != nil {
			u.source = *source

			switch *source {
			case "utm_source=google&utm_medium=cpc":
				u.sourceTag = "Paid Search-Google-Francis-OMS"
			case "utm_source=facebook&utm_medium=paid":
				u.sourceTag = "FB-Neo-OMS"
			case "utm_source=facebook&utm_medium=cpc":
				u.sourceTag = "FB-Weixiong-OMS"
			case "utm_source=facebook&utm_medium=social":
				u.sourceTag = "FB-OMS"
			case "utm_source=instagram&utm_medium=social":
				u.sourceTag = "Ins-OMS"
			case "utm_source=twitter&utm_medium=social":
				u.sourceTag = "Twitter-OMS"
			case "utm_source=youtube&utm_medium=social":
				u.sourceTag = "YouTube-OMS"
			case "utm_source=tiktok&utm_medium=social":
				u.sourceTag = "Tiktok-OMS"
			case "utm_source=hubspot&utm_medium=email":
				u.sourceTag = "EDM-OMS"
			case "Wix":
				// temporary solution for Wix App Review
				u.GuideInfo().SetFinished(true)
			}
		}
	}
}

func WithSourceTag(sourceTag string) Option {
	return func(u *User) {
		u.sourceTag = sourceTag
	}
}

func WithLastLoggedTime(lastLoggedTime time.Time) Option {
	return func(u *User) {
		u.lastLoggedTime = lastLoggedTime
	}
}

func (u *User) SetStatus(status int8) {
	u.status = status
}

func (u *User) SetLastLoggedTime(lastLoggedTime time.Time) {
	u.lastLoggedTime = lastLoggedTime
}

func (u *User) SetStoreCode(storeCode string) {
	u.storeCode = storeCode
}

func (u *User) SetTenant(tenant *domainTenant.Tenant) {
	u.tenant = tenant
}

func (u User) UserID() int64 {
	return u.userID
}

func (u User) Username() string {
	return u.username
}

func (u User) Email() string {
	return u.email
}

func (u User) Password() string {
	return u.password
}

func (u User) Phone() string {
	return u.phone
}

func (u User) Source() string {
	return u.source
}

func (u User) SourceTag() string {
	return u.sourceTag
}

func (u User) Status() int8 {
	return u.status
}

func (u User) LastLoggedTime() time.Time {
	return u.lastLoggedTime
}

func (u User) CreatedAt() time.Time {
	return u.createdAt
}

func (u User) UpdatedAt() time.Time {
	return u.updatedAt
}

func (u User) Tenant() *domainTenant.Tenant {
	return u.tenant
}

func NewUser(
	username string,
	email string,
	password string,
	phone string,
	tenant *domainTenant.Tenant,
	options ...Option,
) (*User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}

	if email == "" {
		return nil, errors.New("email is required")
	}

	if password == "" {
		return nil, errors.New("password is required")
	}

	if phone == "" {
		return nil, errors.New("phone is required")
	}

	if tenant == nil {
		return nil, errors.New("tenant is required")
	}

	user := &User{
		username: username,
		email:    email,
		password: password,
		phone:    phone,
		tenant:   tenant,
		guideInfo: NewGuideInfo(
			1, false, "", []*Question{},
		),
	}

	for _, option := range options {
		option(user)
	}

	return user, nil
}

func UnmarshalUserFromDatabase(
	userID int64,
	username string,
	email string,
	password string,
	phone string,
	tenant *domainTenant.Tenant,
	option ...Option,
) (*User, error) {
	user, err := NewUser(username, email, password, phone, tenant, option...)
	if err != nil {
		return nil, err
	}

	user.userID = userID

	return user, nil
}
