package packer

import (
	"fmt"

	"time"

	slugmaker "github.com/gosimple/slug"
	"github.com/imega-teleport/db2file/indexer"
	"github.com/imega-teleport/db2file/storage"
	"github.com/imega-teleport/db2file/teleport"
	"github.com/imega-teleport/db2file/writer"
	"github.com/yvasiyarov/php_session_decoder/php_serialize"
	"gopkg.in/Masterminds/squirrel.v1"
)

// Packer is interface
type Packer interface {
	Listen(in <-chan interface{}, e chan<- error)
	SaveToFile() error
	SecondSaveToFile() error
	ThirdPackSaveToFile(latest bool) error
}

// Options настройки пакета
type Options struct {
	MaxBytes        int
	PrefixFileName  string
	PathToSave      string
	PathToImages    string
	PrefixTableName string
}

// OptionsExport набор специфических свойств, которые требуется с указания пользователя разместить
type OptionsExport struct {
	Weight    string `json: "weight"`
	Length    string `json: "length"`
	Width     string `json: "width"`
	Height    string `json: "height"`
	TypePrice string `json: "type_price"`
}

type pkg struct {
	Options              Options
	OptionsExport        *OptionsExport
	FirstPack            teleport.FirstPackage
	SecondPack           teleport.SecondPackage
	ThirdPack            teleport.ThirdPackage
	PropertiesCollection propertiesCollection
	Indexer              indexer.Indexer
	FirstPackQty         int
	SecondPackQty        int
	ThirdPackQty         int
	Content              string
}

type propertiesCollection struct {
	ProductID string
	Items     []storage.ProductsProperties
}

// New instance packer
func New(opt Options, optsEx *OptionsExport) Packer {
	return &pkg{
		Options:       opt,
		OptionsExport: optsEx,
		Indexer:       indexer.NewIndexer(),
		FirstPackQty:  1,
		SecondPackQty: 1,
		ThirdPackQty:  1,
	}
}

func (p *pkg) Listen(in <-chan interface{}, e chan<- error) {
	postmeta := map[string]string{
		p.OptionsExport.Length: "_length",
		p.OptionsExport.Height: "_height",
		p.OptionsExport.Weight: "_weight",
		p.OptionsExport.Width:  "_width",
	}
	for v := range in {
		if p.IsFull(p.FirstPack) {
			p.SaveToFile()
			pack := teleport.FirstPackage{}
			p.Content = ""
			p.FirstPack = pack
			p.FirstPackQty++
		}

		if p.SecondIsFull(p.SecondPack) {
			p.SecondSaveToFile()
			pack := teleport.SecondPackage{}
			p.SecondPack = pack
			p.SecondPackQty++
		}

		if p.ThirdPackIsFull(p.ThirdPack) {
			p.ThirdPackSaveToFile(false)
			pack := teleport.ThirdPackage{}
			p.ThirdPack = pack
			p.ThirdPackQty++
		}

		switch v.(type) {
		case storage.Product:
			p.Indexer.Set(teleport.UUID(v.(storage.Product).ID).String())
			p.FirstPack.AddItem(teleport.Post{
				ID:       teleport.UUID(v.(storage.Product).ID),
				AuthorID: 1,
				Date:     time.Now(),
				Content:  v.(storage.Product).Description,
				Title:    v.(storage.Product).Name,
				Excerpt:  "",
				Status:   "publish",
				Name:     v.(storage.Product).Name,
				Modified: time.Now(),
				Type:     "product",
			})
			p.FirstPack.AddItem(teleport.TeleportItem{
				GUID: teleport.UUID(v.(storage.Product).ID),
				Type: "post",
				Date: time.Now(),
			})

			if v.(storage.Product).Article != "" {
				p.ThirdPack.AddItem(teleport.PostMeta{
					PostID: teleport.UUID(v.(storage.Product).ID),
					Key:    "_sku",
					Value:  v.(storage.Product).Article,
				})
			}
			p.ThirdPack.AddItem(teleport.PostMeta{
				PostID: teleport.UUID(v.(storage.Product).ID),
				Key:    "_visibility",
				Value:  "visible",
			})

		case storage.Group:
			p.Indexer.Set(teleport.UUID(v.(storage.Group).ID).String())
			p.FirstPack.AddItem(teleport.Term{
				ID:    teleport.UUID(v.(storage.Group).ID),
				Name:  v.(storage.Group).Name,
				Slug:  teleport.Slug(v.(storage.Group).Name),
				Group: "0",
			})
			p.FirstPack.AddItem(teleport.TeleportItem{
				GUID: teleport.UUID(v.(storage.Group).ID),
				Type: "term",
				Date: time.Now(),
			})
			p.SecondPack.AddItem(teleport.TermTaxonomy{
				TermID:       teleport.UUID(v.(storage.Group).ID),
				Taxonomy:     "product_cat",
				Description:  v.(storage.Group).Name,
				ParentTermID: teleport.UUID(v.(storage.Group).ParentID),
			})

		case storage.ProductsGroups:
			p.ThirdPack.AddItem(teleport.TermRelationship{
				ObjectID:       teleport.UUID(v.(storage.ProductsGroups).ProductID),
				TermTaxonomyID: teleport.UUID(v.(storage.ProductsGroups).GroupID),
			})

		case storage.ProductsProperties:
			if p.PropertiesCollection.ProductID != "" && p.PropertiesCollection.ProductID == v.(storage.ProductsProperties).ProductID {
				p.PropertiesCollection.Items = append(p.PropertiesCollection.Items, v.(storage.ProductsProperties))
			} else {
				if p.PropertiesCollection.ProductID != v.(storage.ProductsProperties).ProductID {
					attrs, _ := p.SerializationProperties(p.PropertiesCollection.Items)
					p.ThirdPack.AddItem(teleport.PostMeta{
						PostID: teleport.UUID(v.(storage.ProductsProperties).ProductID),
						Key:    "_product_attributes",
						Value:  attrs,
					})
				}
				p.PropertiesCollection = propertiesCollection{
					ProductID: v.(storage.ProductsProperties).ProductID,
					Items: []storage.ProductsProperties{
						v.(storage.ProductsProperties),
					},
				}
			}

		case storage.ProductsPropertiesSpecial:
			p.ThirdPack.AddItem(teleport.PostMeta{
				PostID: teleport.UUID(v.(storage.ProductsPropertiesSpecial).ProductID),
				Key:    postmeta[v.(storage.ProductsPropertiesSpecial).PropertyID],
				Value:  v.(storage.ProductsPropertiesSpecial).Value,
			})

		case storage.ProductsPrices:
			p.ThirdPack.AddItem(teleport.PostMeta{
				PostID: teleport.UUID(v.(storage.ProductsPrices).ProductID),
				Key:    "_price",
				Value:  v.(storage.ProductsPrices).Value,
			})
			p.ThirdPack.AddItem(teleport.PostMeta{
				PostID: teleport.UUID(v.(storage.ProductsPrices).ProductID),
				Key:    "_regular_price",
				Value:  v.(storage.ProductsPrices).Value,
			})

		case storage.ProductImage:
			p.productImagePack(v.(storage.ProductImage))
		}
	}
}

func (p *pkg) IsFull(pack teleport.FirstPackage) bool {
	return pack.Length >= p.Options.MaxBytes+p.Indexer.GetLength()+2000
}

func (p *pkg) SecondIsFull(pack teleport.SecondPackage) bool {
	return pack.Length >= p.Options.MaxBytes+2000
}

func (p *pkg) ThirdPackIsFull(pack teleport.ThirdPackage) bool {
	return pack.Length >= p.Options.MaxBytes+2000
}

func (p *pkg) PreContent(s string) {
	p.Content = fmt.Sprintf("%s;", s) + p.Content
}

func (p *pkg) AddContent(s string) {
	p.Content = p.Content + fmt.Sprintf("%s;", s)
}

func (p *pkg) ClearContent() {
	p.Content = ""
}

func (p *pkg) SaveToFile() error {
	p.ClearContent()
	w := writer.NewWriter(p.Options.PrefixFileName, p.Options.PathToSave)
	fileName := w.GetFileName(p.FirstPackQty)
	wpwc := teleport.Wpwc{
		Prefix: p.Options.PrefixTableName,
	}

	idx := indexer.NewIndexer()

	if len(p.FirstPack.Term) > 0 {
		builder := wpwc.BuilderTerm()
		for _, v := range p.FirstPack.Term {
			idx.Set(v.ID.String())
			builder.AddTerm(v)
		}
		p.AddContent(squirrel.DebugSqlizer(builder))
	}

	if len(p.FirstPack.Post) > 0 {
		builder := wpwc.BuilderPost()
		for _, v := range p.FirstPack.Post {
			idx.Set(v.ID.String())
			builder.AddPost(v)
		}
		p.AddContent(squirrel.DebugSqlizer(builder))
	}

	if len(p.Indexer.GetAll()) > 0 {
		builder := wpwc.BuilderTeleportItem()
		for _, v := range p.FirstPack.TeleportItem {
			idx.Set(v.GUID.String())
			builder.AddTeleportItem(v)
		}
		p.AddContent(squirrel.DebugSqlizer(builder))
	}

	if len(idx.GetAll()) > 0 {
		for k := range idx.GetAll() {
			if k != "" {
				p.PreContent(fmt.Sprintf("set @%s=%d", k, p.Indexer.Get(k)))
			}
		}
	}

	p.AddContent("commit")

	p.PreContent(fmt.Sprintf("set @author_id=%d", 1)) //todo author
	p.PreContent(fmt.Sprintf("set @max_post_id=(select ifnull(max(id),0)from %sposts)", p.Options.PrefixTableName))
	p.PreContent(fmt.Sprintf("set @max_term_taxonomy_id=(select ifnull(max(term_taxonomy_id),0)from %sterm_taxonomy)", p.Options.PrefixTableName))
	p.PreContent(fmt.Sprintf("set @max_term_id=(select ifnull(max(term_id),0)from %sterms)", p.Options.PrefixTableName))
	p.PreContent("start transaction")

	if p.FirstPackQty == 1 {
		p.PreContent(fmt.Sprintf("create table if not exists %steleport_item(guid char(80)not null,type char(8)not null,id bigint,date datetime,primary key(`guid`))engine=innodb default charset=utf8", p.Options.PrefixTableName))
	}

	err := w.WriteFile(fileName, p.Content)
	return err
}

func (p *pkg) SecondSaveToFile() error {
	p.ClearContent()
	w := writer.NewWriter(fmt.Sprintf("sec/%s", p.Options.PrefixFileName), p.Options.PathToSave)
	fileName := w.GetFileName(p.SecondPackQty)
	wpwc := teleport.Wpwc{
		Prefix: p.Options.PrefixTableName,
	}

	idx := indexer.NewIndexer()

	if len(p.SecondPack.TermTaxonomy) > 0 {
		builder := wpwc.BuilderTermTaxonomy()
		for _, v := range p.SecondPack.TermTaxonomy {
			idx.Set(v.TermID.String())
			idx.Set(v.ParentTermID.String())
			builder.AddTermTaxonomy(v)
		}
		p.AddContent(squirrel.DebugSqlizer(builder))
	}

	if len(idx.GetAll()) > 0 {
		for k := range idx.GetAll() {
			if k != "" {
				p.PreContent(fmt.Sprintf("set @%s=(select id from %steleport_item where guid='%s')", k, wpwc.Prefix, k))
			}
		}
	}

	err := w.WriteFile(fileName, p.Content)
	return err
}

func (p *pkg) ThirdPackSaveToFile(latest bool) error {
	p.ClearContent()
	p.ThirdPackToContent(latest)
	err := p.Write("thi", p.ThirdPackQty)
	return err
}

// ThirdPackToContent подготовка содержимого пакета к записи в файл
func (p *pkg) ThirdPackToContent(latest bool) error {
	wpwc := teleport.Wpwc{
		Prefix: p.Options.PrefixTableName,
	}

	if latest && len(p.PropertiesCollection.Items) > 0 {
		attrs, _ := p.SerializationProperties(p.PropertiesCollection.Items)
		p.ThirdPack.AddItem(teleport.PostMeta{
			PostID: teleport.UUID(p.PropertiesCollection.ProductID),
			Key:    "_product_attributes",
			Value:  attrs,
		})
	}

	idx := indexer.NewIndexer()
	idxTermTaxonomy := indexer.NewIndexer()
	idxPost := indexer.NewIndexer()

	if len(p.ThirdPack.TermRelationship) > 0 {
		builder := wpwc.BuilderTermRelationships()
		for _, v := range p.ThirdPack.TermRelationship {
			idxPost.Set(v.ObjectID.String())
			idxTermTaxonomy.Set(v.TermTaxonomyID.String())
			builder.AddTermRelationships(v)
		}
		p.AddContent(squirrel.DebugSqlizer(builder))
	}

	if len(p.ThirdPack.PostMeta) > 0 {
		builder := wpwc.BuilderPostMeta()
		for _, v := range p.ThirdPack.PostMeta {
			idxPost.Set(v.PostID.String())
			builder.AddrPostMeta(v)
		}
		p.AddContent(squirrel.DebugSqlizer(builder))
	}

	if len(idxTermTaxonomy.GetAll()) > 0 {
		for k := range idxTermTaxonomy.GetAll() {
			if k != "" {
				p.PreContent(fmt.Sprintf("set @%s=(select term_taxonomy_id from wp_term_taxonomy where term_id=(select id from %steleport_item where guid='%s'))", k, wpwc.Prefix, k))
			}
		}
	}

	if len(idxPost.GetAll()) > 0 {
		for k := range idxPost.GetAll() {
			if k != "" {
				p.PreContent(fmt.Sprintf("set @%s=(select id from %steleport_item where guid='%s')", k, wpwc.Prefix, k))
			}
		}
	}

	if len(p.ThirdPack.Post) > 0 {
		builder := wpwc.BuilderPost()
		for _, v := range p.ThirdPack.Post {
			idx.Set(v.ID.String())
			builder.AddPost(v)
		}
		p.PreContent(squirrel.DebugSqlizer(builder))
	}

	if len(idx.GetAll()) > 0 {
		for k := range idx.GetAll() {
			if k != "" {
				p.PreContent(fmt.Sprintf("set @%s=%d", k, idx.Get(k)))
			}
		}
	}

	p.PreContent(fmt.Sprintf("set @max_post_id=(select ifnull(max(id),0)from %sposts)", p.Options.PrefixTableName))

	return nil
}

func (p *pkg) Write(name string, idx int) error {
	w := writer.NewWriter(fmt.Sprintf("%s/%s", name, p.Options.PrefixFileName), p.Options.PathToSave)
	fileName := w.GetFileName(idx)
	err := w.WriteFile(fileName, p.Content)
	return err
}

//@todo в teleport
func (p *pkg) SerializationProperties(items []storage.ProductsProperties) (string, error) {
	encoder := php_serialize.NewSerializer()
	source := map[php_serialize.PhpValue]php_serialize.PhpValue{}

	for _, v := range items {
		source[slugmaker.Make(v.PropertyName)] = map[php_serialize.PhpValue]php_serialize.PhpValue{
			"name":         v.PropertyName,
			"value":        v.Value,
			"position":     "0",
			"is_visible":   "1",
			"is_variation": "0",
			"is_taxonomy":  "0",
		}
	}

	return encoder.Encode(source)
}
