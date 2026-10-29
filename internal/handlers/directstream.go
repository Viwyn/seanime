package handlers

import (
	"fmt"
	"net/http"
	"seanime/internal/database/db_bridge"
	"seanime/internal/directstream"
	"seanime/internal/mediastream/videofile"
	"seanime/internal/mkvparser"
	"seanime/internal/security"
	"seanime/internal/util"

	"github.com/labstack/echo/v4"
)

// HandleDirectstreamPlayLocalFile
//
//	@summary request local file stream.
//	@desc This requests a local file stream and returns the media container to start the playback.
//	@returns mediastream.MediaContainer
//	@route /api/v1/directstream/play/localfile [POST]
func (h *Handler) HandleDirectstreamPlayLocalFile(c echo.Context) error {
	if err := h.guardMediaConsumption(c); err != nil {
		return err
	}

	type body struct {
		Path     string `json:"path"`     // The path of the file.
		ClientId string `json:"clientId"` // The session id
	}

	var b body
	if err := c.Bind(&b); err != nil {
		return h.RespondWithError(c, err)
	}

	b.ClientId = getRequestClientId(c, b.ClientId)
	b.Path = util.ResolvePhysicalPath(b.Path)
	if err := h.guardStrictFilesystemPath(c, b.Path); err != nil {
		return err
	}

	lfs, _, err := db_bridge.GetLocalFiles(h.App.Database)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.App.DirectStreamManager.PlayLocalFile(c.Request().Context(), directstream.PlayLocalFileOptions{
		ClientId:   b.ClientId,
		Path:       b.Path,
		LocalFiles: lfs,
	})
}

// HandleDirectstreamTranscode
//
//	@summary transcode the current torrent stream.
//	@desc This starts transcoding the torrent currently being streamed so browsers that can't play its codecs can watch it.
//	@desc Subtitles keep coming from the direct stream, only the video and audio are transcoded.
//	@returns mediastream.MediaContainer
//	@route /api/v1/directstream/transcode [POST]
func (h *Handler) HandleDirectstreamTranscode(c echo.Context) error {
	if err := h.guardMediaConsumption(c); err != nil {
		return err
	}

	type body struct {
		ClientId string `json:"clientId"` // The session id
	}

	var b body
	if err := c.Bind(&b); err != nil {
		return h.RespondWithError(c, err)
	}

	b.ClientId = getRequestClientId(c, b.ClientId)

	if err := h.guardPrivilegedMediastream(c, h.App.SecondarySettings.Mediastream); err != nil {
		return err
	}

	inputUrl, id, ok := h.App.TorrentstreamRepository.GetTranscodeInput()
	if !ok {
		return h.RespondWithError(c, fmt.Errorf("no torrent is being streamed"))
	}

	mediaContainer, err := h.App.MediastreamRepository.RequestTranscodeRemoteStream(inputUrl, videofile.GetHashForStream(id), b.ClientId)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, mediaContainer)
}

// HandleDirectstreamConvertSubs
//
//	@summary converts subtitles from one format to another.
//	@returns string
//	@route /api/v1/directstream/subs/convert-subs [POST]
func (h *Handler) HandleDirectstreamConvertSubs(c echo.Context) error {
	type body struct {
		Url     string `json:"url"`
		Content string `json:"content"`
		To      string `json:"to"`
	}

	var b body
	if err := c.Bind(&b); err != nil {
		return h.RespondWithError(c, err)
	}

	if b.Url == "" && b.Content == "" {
		return h.RespondWithError(c, fmt.Errorf("url or content is required"))
	}

	if b.To == "" {
		return h.RespondWithError(c, fmt.Errorf("to is required"))
	}

	to := mkvparser.SubtitleTypeASS
	switch b.To {
	case "ass":
		to = mkvparser.SubtitleTypeASS
	case "vtt":
		to = mkvparser.SubtitleTypeWEBVTT
	}

	if len(b.Content) > 0 {
		// Convert from content
		ret, err := h.App.VideoCore.ConvertSubsTo(b.Content, mkvparser.SubtitleTypeUnknown, to)
		if err != nil {
			return h.RespondWithError(c, err)
		}
		return h.RespondWithData(c, ret)
	}

	// Convert from url
	if err := security.ValidateOutboundUrl(b.Url); err != nil {
		return h.RespondWithStatusError(c, echo.ErrForbidden.Code, err)
	}

	ret, err := h.App.VideoCore.FetchAndConvertSubsTo(b.Url, to)
	if err != nil {
		return h.RespondWithError(c, err)
	}

	return h.RespondWithData(c, ret)
}

func (h *Handler) HandleDirectstreamGetStream() http.Handler {
	streamHandler := h.App.DirectStreamManager.ServeEchoStream()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h != nil && h.App != nil && h.App.Config != nil && !canConsumeMedia(r, h.App.Config.Server.Password, h.App.Config.Server.AccessAllowlist) {
			http.Error(w, errPrivilegedExecutionDenied.Error(), http.StatusForbidden)
			return
		}

		streamHandler.ServeHTTP(w, r)
	})
}

func (h *Handler) HandleDirectstreamGetAttachments(c echo.Context) error {
	if err := h.guardMediaConsumption(c); err != nil {
		return err
	}

	return h.App.DirectStreamManager.ServeEchoAttachments(c)
}
