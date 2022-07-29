package ydatabase

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"text/template"
	"time"

	environ "github.com/ydb-platform/ydb-go-sdk-auth-environ"
	ydb "github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

type YandexDatabase struct {
	ctx         context.Context
	prefix      string
	tableclient table.Client
	cancel      context.CancelFunc
	db          ydb.Connection
}

type YDBTableMessage struct {
	msg_id uint32
	board  [16]uint8
	time   time.Time
}

func NewYandexDatabase(endpoint string) YandexDatabase {
	ctx, cancel := context.WithCancel(context.Background())

	db, err := ydb.Open(
		ctx,
		endpoint,
		environ.WithEnvironCredentials(ctx),
	)

	if err != nil {
		panic(fmt.Errorf("connect error: %w", err))
	}

	prefix := path.Join(db.Name(), "tg_tag_bot")

	return YandexDatabase{
		ctx:         ctx,
		prefix:      prefix,
		tableclient: db.Table(),
		cancel:      cancel,
		db:          db,
	}
}

func (yd *YandexDatabase) create_database(name string) error {
	return yd.tableclient.Do(
		yd.ctx,
		func(ctx context.Context, s table.Session) error {
			return s.CreateTable(ctx, path.Join(yd.prefix, name),
				options.WithColumn("msg_id", types.Optional(types.TypeUint32)),
				options.WithColumn("P1", types.Optional(types.TypeUint8)),
				options.WithColumn("P2", types.Optional(types.TypeUint8)),
				options.WithColumn("P3", types.Optional(types.TypeUint8)),
				options.WithColumn("P4", types.Optional(types.TypeUint8)),
				options.WithColumn("P5", types.Optional(types.TypeUint8)),
				options.WithColumn("P6", types.Optional(types.TypeUint8)),
				options.WithColumn("P7", types.Optional(types.TypeUint8)),
				options.WithColumn("P8", types.Optional(types.TypeUint8)),
				options.WithColumn("P9", types.Optional(types.TypeUint8)),
				options.WithColumn("P10", types.Optional(types.TypeUint8)),
				options.WithColumn("P11", types.Optional(types.TypeUint8)),
				options.WithColumn("P12", types.Optional(types.TypeUint8)),
				options.WithColumn("P13", types.Optional(types.TypeUint8)),
				options.WithColumn("P14", types.Optional(types.TypeUint8)),
				options.WithColumn("P15", types.Optional(types.TypeUint8)),
				options.WithColumn("P16", types.Optional(types.TypeUint8)),
				options.WithColumn("time", types.Optional(types.TypeDatetime)),
				options.WithPrimaryKeyColumn("msg_id"),
			)
		},
	)
}

func (yd *YandexDatabase) write_to_database(name string, data [1]YDBTableMessage) error {
	var wrap = func(err error, explanation string) error {
		if err != nil {
			return fmt.Errorf("%s: %w", explanation, err)
		}
		return err
	}

	err := yd.tableclient.Do(
		yd.ctx,
		func(ctx context.Context, session table.Session) error {
			rows := make([]types.Value, 0, len(data))

			for _, msg := range data {

				rows = append(rows, types.StructValue(
					types.StructFieldValue("msg_id", types.Uint32Value(uint32(msg.msg_id))),
					types.StructFieldValue("P1", types.Uint8Value(uint8(msg.board[0]))),
					types.StructFieldValue("P2", types.Uint8Value(uint8(msg.board[1]))),
					types.StructFieldValue("P3", types.Uint8Value(uint8(msg.board[2]))),
					types.StructFieldValue("P4", types.Uint8Value(uint8(msg.board[3]))),
					types.StructFieldValue("P5", types.Uint8Value(uint8(msg.board[4]))),
					types.StructFieldValue("P6", types.Uint8Value(uint8(msg.board[5]))),
					types.StructFieldValue("P7", types.Uint8Value(uint8(msg.board[6]))),
					types.StructFieldValue("P8", types.Uint8Value(uint8(msg.board[7]))),
					types.StructFieldValue("P9", types.Uint8Value(uint8(msg.board[8]))),
					types.StructFieldValue("P10", types.Uint8Value(uint8(msg.board[9]))),
					types.StructFieldValue("P11", types.Uint8Value(uint8(msg.board[10]))),
					types.StructFieldValue("P12", types.Uint8Value(uint8(msg.board[11]))),
					types.StructFieldValue("P13", types.Uint8Value(uint8(msg.board[12]))),
					types.StructFieldValue("P14", types.Uint8Value(uint8(msg.board[13]))),
					types.StructFieldValue("P15", types.Uint8Value(uint8(msg.board[14]))),
					types.StructFieldValue("P16", types.Uint8Value(uint8(msg.board[15]))),
					types.StructFieldValue("time", types.DatetimeValueFromTime(msg.time)),
				))
			}

			return wrap(session.BulkUpsert(ctx, path.Join(yd.prefix, name), types.ListValue(rows...)),
				"failed to perform bulk upsert")
		})
	return wrap(err, "failed to write log batch")
}

func (yd *YandexDatabase) read_from_table(name string, msg_id uint32, data *YDBTableMessage) (err error) {
	type templateConfig struct {
		TablePathPrefix string
	}
	var render = func(t *template.Template, data interface{}) string {
		var buf bytes.Buffer
		err := t.Execute(&buf, data)
		if err != nil {
			panic(err)
		}
		return buf.String()
	}
	query := render(
		template.Must(template.New("").Parse(
			"PRAGMA TablePathPrefix('{{ .TablePathPrefix}}');\n"+
				"DECLARE $msg_id AS Uint32; "+
				"SELECT * "+
				"FROM `HALLO_BASE`\n"+
				"WHERE msg_id == $msg_id",
		)),
		templateConfig{
			TablePathPrefix: yd.prefix,
		},
	)

	readTx := table.TxControl(
		table.BeginTx(
			table.WithOnlineReadOnly(),
		),
		table.CommitTx(),
	)

	var res result.Result
	err = yd.tableclient.Do(
		yd.ctx,
		func(ctx context.Context, s table.Session) (err error) {
			_, res, err = s.Execute(ctx, readTx, query,
				table.NewQueryParameters(
					table.ValueParam("$msg_id", types.Uint32Value(msg_id)),
				),
				options.WithCollectStatsModeBasic(),
			)
			return
		},
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = res.Close()
	}()

	for res.NextResultSet(yd.ctx) {
		for res.NextRow() {
			err = res.ScanNamed(
				named.Required("msg_id", &data.msg_id),
				named.Required("P1", &data.board[0]),
				named.Required("P2", &data.board[1]),
				named.Required("P3", &data.board[2]),
				named.Required("P4", &data.board[3]),
				named.Required("P5", &data.board[4]),
				named.Required("P6", &data.board[5]),
				named.Required("P7", &data.board[6]),
				named.Required("P8", &data.board[7]),
				named.Required("P9", &data.board[8]),
				named.Required("P10", &data.board[9]),
				named.Required("P11", &data.board[10]),
				named.Required("P12", &data.board[11]),
				named.Required("P13", &data.board[12]),
				named.Required("P14", &data.board[13]),
				named.Required("P15", &data.board[14]),
				named.Required("P16", &data.board[15]),
			)
			if err != nil {
				panic(err)
			}

		}
	}
	return res.Err()
}
