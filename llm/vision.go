package llm

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"net/http"

	"github.com/openai/openai-go/v3"
	"golang.org/x/image/webp"
)

const captionPrompt = "Describe this image in 2-3 sentences for a search index. Focus on the visible subject - what it is, key details, style. Do not interpret or speculate; describe only what is shown."

func (c *Client) HasVision() bool {
	return c.cfg.VisionModel != ""
}

func convertWebPToPNG(webpBytes []byte) ([]byte, error) {
	img, err := webp.Decode(bytes.NewReader(webpBytes))
	if err != nil {
		return nil, fmt.Errorf("decode webp: %w", err)
	}

	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("encode png: %w", err)
	}

	return buf.Bytes(), nil
}

func (c *Client) DescribeImage(ctx context.Context, mime string, image []byte) (string, error) {
	if c.cfg.VisionModel == "" {
		return "", errors.New("no vision model configured")
	}

	if len(image) == 0 {
		return "", errors.New("empty image")
	}

	if mime == "" {
		mime = http.DetectContentType(image)
	}

	// If it's WebP, convert to PNG
	if mime == "image/webp" {
		pngBytes, err := convertWebPToPNG(image)
		if err != nil {
			return "", fmt.Errorf("webp → png conversion failed: %w", err)
		}
		image = pngBytes
		mime = "image/png"
	}

	dataURL := "data:" + mime + ";base64," + base64.StdEncoding.EncodeToString(image)

	resp, err := c.sdk.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: c.cfg.VisionModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage([]openai.ChatCompletionContentPartUnionParam{
				openai.TextContentPart(captionPrompt),
				openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
					URL: dataURL,
				}),
			}),
		},
	})
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices found")
	}

	return resp.Choices[0].Message.Content, nil
}
