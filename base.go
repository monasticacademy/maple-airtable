package airtable

import (
	"context"
	"encoding/json"
	"net/url"
	"os"
)

// Base type of airtable base.
type Base struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	PermissionLevel string `json:"permissionLevel"`
}

// Base type of airtable bases.
type Bases struct {
	Bases  []*Base `json:"bases"`
	Offset string  `json:"offset,omitempty"`
}

type Field struct {
	ID          string         `json:"id,omitempty"`
	Type        string         `json:"type"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Options     map[string]any `json:"options,omitempty"`
}

type View struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type TableSchema struct {
	ID             string   `json:"id,omitempty"`
	PrimaryFieldID string   `json:"primaryFieldId,omitempty"`
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	Fields         []*Field `json:"fields"`
	Views          []*View  `json:"views,omitempty"`
}

type Tables struct {
	Tables []*TableSchema `json:"tables"`
}

// GetBasesWithParams get bases with url values params
// https://airtable.com/developers/web/api/list-bases
func (at *Client) GetBasesWithParams(params url.Values) (*Bases, error) {
	return at.GetBasesWithParamsContext(context.Background(), params)
}

// getBasesWithParamsContext get bases with url values params
// with custom context
func (at *Client) GetBasesWithParamsContext(ctx context.Context, params url.Values) (*Bases, error) {
	bases := new(Bases)

	err := at.get(ctx, "meta", "bases", "", params, bases)
	if err != nil {
		return nil, err
	}

	return bases, nil
}

// Table represents table object.
type BaseConfig struct {
	client *Client
	dbId   string
}

// GetBase return Base object.
func (c *Client) GetBaseSchema(dbId string) *BaseConfig {
	return &BaseConfig{
		client: c,
		dbId:   dbId,
	}
}

// Do send the prepared
func (b *BaseConfig) Do() (*Tables, error) {
	return b.GetTables()
}

// Do send the prepared with custom context
func (b *BaseConfig) DoContext(ctx context.Context) (*Tables, error) {
	return b.GetTablesContext(ctx)
}

// GetTables get tables from a base with url values params
// https://airtable.com/developers/web/api/get-base-schema
func (b *BaseConfig) GetTables() (*Tables, error) {
	return b.GetTablesContext(context.Background())
}

// getTablesContext get tables from a base with url values params
// with custom context
func (b *BaseConfig) GetTablesContext(ctx context.Context) (*Tables, error) {
	tables := new(Tables)

	err := b.client.get(ctx, "meta/bases", b.dbId, "tables", nil, tables)
	if err != nil {
		return nil, err
	}

	return tables, nil
}

type CreateBaseRequest struct {
	Name      string         `json:"name"`
	Workspace string         `json:"workspaceId"`
	Tables    []*TableSchema `json:"tables"`
}

type CreateBaseResponse struct {
	BaseID string         `json:"id"` // id of created base
	Tables []*TableSchema `json:"tables"`
}

// CreateBases create base from table schemas
// https://airtable.com/developers/web/api/create-base
func (at *Client) CreateBase(req *CreateBaseRequest) (*CreateBaseResponse, error) {
	return at.CreateBaseContext(context.Background(), req)
}

// CreateBasesContext create base from table schema with context
// https://airtable.com/developers/web/api/create-base
func (at *Client) CreateBaseContext(ctx context.Context, req *CreateBaseRequest) (*CreateBaseResponse, error) {
	resp := new(CreateBaseResponse)

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(req)

	err := at.post(ctx, "meta", "bases", req, resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
