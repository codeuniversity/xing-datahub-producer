package handler

import (
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/codeuniversity/xing-datahub-protocol"
	"github.com/golang/protobuf/jsonpb"

	"net/http/httptest"

	"github.com/Shopify/sarama"
	"github.com/Shopify/sarama/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRequestHandler(t *testing.T) {
	Convey("Given a wrong user", t, func() {
		user := &protocol.User{Id: 1}
		marshaler := &jsonpb.Marshaler{}
		message, err := marshaler.MarshalToString(user)
		So(err, ShouldBeNil)
		reader := strings.NewReader(message)
		config := sarama.NewConfig()
		config.Producer.Return.Successes = false
		config.Producer.Return.Errors = false
		producer := mocks.NewAsyncProducer(t, config)
		handler := RequestHandler{
			Producer:   producer,
			RawMessage: &protocol.RawUser{},
			Topic:      "users",
		}
		Convey("With no access-token", func() {
			req := httptest.NewRequest("POST", "/users", reader)
			resp := httptest.NewRecorder()
			Convey("When there is a token defined in os.env", func() {
				os.Setenv("token", "123")
				handler.ServeHTTP(resp, req)

				So(resp.Code, ShouldEqual, http.StatusUnauthorized)
			})

			Convey("When there is no token defined in os.env", func() {
				os.Setenv("token", "")
				handler.ServeHTTP(resp, req)

				So(resp.Code, ShouldNotEqual, http.StatusUnauthorized)

			})
		})
		Convey("With the correct access-token", func() {
			req := httptest.NewRequest("POST", "/users", reader)
			req.Header.Add("access-token", "123")
			resp := httptest.NewRecorder()
			Convey("When there is a token defined in os.env", func() {
				os.Setenv("token", "123")
				handler.ServeHTTP(resp, req)

				So(resp.Code, ShouldNotEqual, http.StatusUnauthorized)
			})

			Convey("When there is no token defined in os.env", func() {
				os.Setenv("token", "")
				handler.ServeHTTP(resp, req)

				So(resp.Code, ShouldNotEqual, http.StatusUnauthorized)
			})
			Convey("We answer with the correct code", func() {
				handler.ServeHTTP(resp, req)
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}
