package packer

import (
	"fmt"
	"testing"
	"time"

	"github.com/imega-teleport/db2file/indexer"
	"github.com/imega-teleport/db2file/storage"
	"github.com/stretchr/testify/assert"
)

func Test_PackProductImage(t *testing.T) {
	p := pkg{
		Indexer: indexer.NewIndexer(),
		Options: Options{
			PathToImages: "/go/src/github.com/imega-teleport/db2file",
		},
	}

	item := storage.ProductImage{
		ProductID: "2d551104-b0ef-11e3-91a0-94de8026f172",
		URL:       "imager/teleport.png",
	}

	err := p.productImagePack(item)
	assert.NoError(t, err)

	err = p.ThirdPackToContent(true)
	assert.NoError(t, err)

	expected := fmt.Sprintf("set @max_post_id=(select ifnull(max(id),0)from posts);"+
		"set @2d551104b0ef11e391a094de8026f172_img=1;"+
		"INSERT INTO posts (id,post_author,post_date,post_date_gmt,post_content,post_title,post_excerpt,post_status,post_name,post_modified,post_modified_gmt,post_parent,post_type,post_mime_type) "+
		"VALUES (@max_post_id+@2d551104b0ef11e391a094de8026f172_img,"+
		"1,"+
		"'%s',"+
		"'%s',"+
		"'',"+
		"'/go/src/github.com/imega-teleport/db2file/imager/teleport.png',"+
		"'',"+
		"'inherit',"+
		"'/go/src/github.com/imega-teleport/db2file/imager/teleport.png',"+
		"'%s',"+
		"'%s',"+
		"@2d551104b0ef11e391a094de8026f172,"+
		"'attachment',"+
		"'image/png');", time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"), time.Now().Format("2006-01-02 15:04:05"))

	assert.Equal(t, expected, p.Content)
}
