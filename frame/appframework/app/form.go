package app

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/zanlichard/beegoe/validation"
)

func BindOnly(c *gin.Context, form interface{}) error {
	if err := c.Bind(form); err != nil {
		return err
	}
	return nil
}

func ValidOnly(c *gin.Context, form interface{}) error {
	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return err
	}
	if !check {
		MarkErrors(c, valid.Errors)
		return BuildFormErr(valid.Errors)
	}
	return nil
}

func BindAndValid(c *gin.Context, form interface{}) error {
	if err := c.Bind(form); err != nil {
		return err
	}
	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return err
	}
	if !check {
		MarkErrors(c, valid.Errors)
		return BuildFormErr(valid.Errors)
	}
	return nil
}

func BuildFormErr(errs []*validation.Error) error {
	var msg strings.Builder
	for _, v := range errs {
		msg.WriteString(v.Error())
	}
	return errors.New(msg.String())
}
