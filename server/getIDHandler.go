package server

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hoshinonyaruko/gensokyo/idmap"
	"github.com/hoshinonyaruko/gensokyo/mylog"
)

func GetIDHandler(c *gin.Context) {
	idOrRow := c.Query("id")
	typeVal, err := strconv.Atoi(c.Query("type"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
		return
	}

	switch typeVal {
	case 1:
		newRow, err := idmap.StoreIDv2(idOrRow)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"row": newRow})

	case 2:
		id, err := idmap.RetrieveRowByIDv2(idOrRow)
		if err == idmap.ErrKeyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "ID not found"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})

	case 3:
		// 存储
		section := c.Query("id")
		subtype := c.Query("subtype")
		value := c.Query("value")
		err := idmap.WriteConfigv2(section, subtype, value)
		if err != nil {
			mylog.Printf(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})

	case 4:
		// 获取值
		section := c.Query("id")
		subtype := c.Query("subtype")
		value, err := idmap.ReadConfigv2(section, subtype)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"value": value})

	case 5:
		oldRowValue, err := strconv.ParseInt(c.Query("oldRowValue"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oldRowValue"})
			return
		}

		newRowValue, err := strconv.ParseInt(c.Query("newRowValue"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid newRowValue"})
			return
		}

		err = idmap.UpdateVirtualValuev2(oldRowValue, newRowValue)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success"})

	case 6:
		virtualValue, err := strconv.ParseInt(c.Query("virtualValue"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid virtualValue"})
			return
		}

		virtual, real, err := idmap.RetrieveRealValuev2(virtualValue)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"virtual": virtual, "real": real})
	case 7:
		realValue := c.Query("id")
		if realValue == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

	}

}
