package middleware

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type request struct {
	A string
	B int
}

type response struct {
	C string
	D int
}

func TestCompress_ReceiveGZIP(t *testing.T) {
	t.Run("Returns gzip encoded response", func(t *testing.T) {
		req := request{"request", 100}
		reqJson, err := encodeRequestJson(req)
		require.NoError(t, err)

		e := echo.New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", reqJson)
		r.Header.Set(echo.HeaderAcceptEncoding, "gzip")
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		ctx := e.NewContext(r, w)
		handler := func(c echo.Context) error {
			var reqInner request
			if err := c.Bind(&reqInner); err != nil {
				t.Log(reqInner)
				t.Error(err)
				return echo.ErrBadRequest
			}
			assert.Equal(t, req, reqInner)
			resp := response{"response", 200}
			return c.JSON(http.StatusOK, resp)
		}
		err = Compression()(handler)(ctx)

		require.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, 200)
		resp, err := decodeGzipResponse(w)
		require.NoError(t, err)
		assert.Equal(t, resp.C, "response")
		assert.Equal(t, resp.D, 200)
	})

	t.Run("Accepts gzip encoded request", func(t *testing.T) {
		req := request{"request", 100}
		gzipRequest, err := encodeRequestGzip(req)
		require.NoError(t, err)

		e := echo.New()
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", gzipRequest)
		r.Header.Set(echo.HeaderContentEncoding, "gzip")
		r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

		ctx := e.NewContext(r, w)
		handler := func(c echo.Context) error {
			var reqInner request
			if err := c.Bind(&reqInner); err != nil {
				t.Log(reqInner)
				t.Error(err)
				return echo.ErrBadRequest
			}
			assert.Equal(t, req, reqInner)
			resp := response{"response", 200}
			return c.JSON(http.StatusOK, resp)
		}
		err = Compression()(handler)(ctx)

		require.NoError(t, err)
		assert.Equal(t, ctx.Response().Status, 200)
		resp, err := decodeJsonResponse(w)
		require.NoError(t, err)
		assert.Equal(t, resp.C, "response")
		assert.Equal(t, resp.D, 200)
	})
}

func decodeGzipResponse(w *httptest.ResponseRecorder) (response, error) {
	var resp response
	var gzipReader *gzip.Reader
	var err error
	if gzipReader, err = gzip.NewReader(w.Result().Body); err != nil {
		return resp, err
	}
	err = json.NewDecoder(gzipReader).Decode(&resp)
	return resp, err
}

func decodeJsonResponse(w *httptest.ResponseRecorder) (response, error) {
	var resp response
	err := json.NewDecoder(w.Result().Body).Decode(&resp)
	return resp, err
}

func encodeRequestJson(req request) (*bytes.Buffer, error) {
	var b []byte
	reqJson := bytes.NewBuffer(b)
	err := json.NewEncoder(reqJson).Encode(req)
	return reqJson, err
}

func encodeRequestGzip(req request) (*bytes.Buffer, error) {
	var b []byte
	var jsonBody []byte
	var err error
	if jsonBody, err = json.Marshal(req); err != nil {
		return nil, err
	}
	reqGzip := bytes.NewBuffer(b)
	gzipWriter := gzip.NewWriter(reqGzip)
	if _, err = gzipWriter.Write(jsonBody); err != nil {
		return nil, err
	}
	gzipWriter.Close()
	return reqGzip, nil
}
