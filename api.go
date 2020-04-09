package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg"
)

// SuperHandler interface for REST API for Super
type SuperHandler interface {
	supersPOSTHandler(c *gin.Context)
	supersGETHandler(c *gin.Context)
	supersPUTHandler(c *gin.Context)
	supersDeleteHandler(c *gin.Context)
}

// SuperAPI implements SuperHandler interface
type SuperAPI struct {
	DB     *pg.DB
	Router *gin.Engine
}

type errorResponseJSON struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func (api *SuperAPI) handleSuperCreate(c *gin.Context, super *Super) {
	if _, err := super.Create(api.DB); err != nil {
		if _, ok := err.(*errorSuperAlreadyExists); ok {
			c.JSON(http.StatusConflict, errorResponseJSON{
				"Super already exists - update it instead",
				err.Error(),
			})
		} else {
			panic(err)
		}
	} else {
		c.JSON(http.StatusCreated, super)
	}
}

func (api *SuperAPI) handleSuperBindingJSON(c *gin.Context) (*Super, bool) {
	super := Super{}

	if err := c.ShouldBindJSON(&super); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseJSON{
			"Error processing the payload",
			err.Error(),
		})
		return &super, false
	}
	return &super, true
}

func (api *SuperAPI) superHeroPOSTHandler(c *gin.Context) {

	super, ok := api.handleSuperBindingJSON(c)
	if !ok {
		return
	}
	super.Type = "HERO"

	api.handleSuperCreate(c, super)
}

func (api *SuperAPI) superVilanPOSTHandler(c *gin.Context) {

	super, ok := api.handleSuperBindingJSON(c)
	if !ok {
		return
	}
	super.Type = "VILAN"

	api.handleSuperCreate(c, super)
}

func (api *SuperAPI) supersPOSTHandler(c *gin.Context) {

	super, ok := api.handleSuperBindingJSON(c)
	if !ok {
		return
	}

	api.handleSuperCreate(c, super)
}

// supersGETHandler handles GET Requests.
// Valid requests:
// 	GET /supers/name
// 	GET /supers/uuid
func (api *SuperAPI) supersGETHandler(c *gin.Context) {

	if c.Param("id") != "" { // Searching for a specific user (by name or uuid)
		var super *Super
		var err error

		super, err = super.getByNameOrUUID(api.DB, c.Param("id"))
		if err != nil {
			if _, ok := err.(*errorSuperNotFound); ok {
				c.JSON(http.StatusNotFound, errorResponseJSON{
					"No Super was found",
					err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, errorResponseJSON{
					"Internal Server Error",
					err.Error(),
				})
			}
		}

		c.JSON(http.StatusOK, super)

	} else {

		sFilter := Super{}

		if err := c.ShouldBind(&sFilter); err != nil {
			c.JSON(http.StatusBadRequest, errorResponseJSON{
				"Could not process Payload (query parameters)",
				err.Error(),
			})
		} else {
			results := sFilter.ReadAll(api.DB)

			c.JSON(http.StatusOK, results)
		}
	}
}

func (api *SuperAPI) supersPUTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not Implemented",
		"Updating Super is out of scope for now",
	})
}

func (api *SuperAPI) supersDeleteHandler(c *gin.Context) {
	super := new(Super)
	err := super.DeleteByNameOrUUID(api.DB, c.Param("id"))

	if err != nil {
		if _, ok := err.(*errorSuperNotFound); ok {
			// nothing was deleted - return 404
			c.JSON(http.StatusNotFound, errorResponseJSON{
				"Super Not Found",
				err.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, errorResponseJSON{
				"Unexpected Error",
				err.Error(),
			})
		}
	} else {
		// Deleted. No Content is needed
		c.Status(http.StatusNoContent)
	}
}

func (api *SuperAPI) groupsPOSTHandler(c *gin.Context) {
	group := Group{}

	if err := c.ShouldBindJSON(&group); err != nil {
		c.JSON(http.StatusBadRequest, errorResponseJSON{
			"Error processing the payload",
			err.Error(),
		})
		return
	}

	if group, err := group.Create(api.DB); err != nil {
		if _, ok := err.(*errorGroupAlreadyExists); ok {
			c.JSON(http.StatusConflict, errorResponseJSON{
				"Group already exists - update it instead",
				err.Error(),
			})
		} else if _, ok := err.(*errorGroupSuperRelation); ok {
			log.Println("Found some non-fatal errors. Will log and ignore:", err.Error())
		} else {
			panic(err)
		}
	} else {
		c.JSON(http.StatusCreated, group)
	}

}

func (api *SuperAPI) groupsGETHandler(c *gin.Context) {
	var group *Group
	var err error

	group, err = group.GetByName(api.DB, c.Param("name"))
	if err != nil {
		if _, ok := err.(*errorGroupNotFound); ok {
			c.JSON(http.StatusNotFound, errorResponseJSON{
				"Group not found", err.Error(),
			})
		} else {
			panic(err)
		}
	} else {
		c.JSON(http.StatusOK, group)
	}
}

func (api *SuperAPI) groupsPUTHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not implemented", "out of the scope for this exercise",
	})
}

func (api *SuperAPI) groupsDeleteHandler(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, errorResponseJSON{
		"Not implemented", "out of the scope for this exercise",
	})
}
