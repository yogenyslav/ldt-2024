package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) processChunk(wg *sync.WaitGroup, c tele.Context, out <-chan chatresp.Response) {
	defer wg.Done()
	buff := strings.Builder{}
	for {
		select {
		case <-time.After(time.Minute):
			if err := c.Send("response timeout exceeded"); err != nil {
				log.Error().Err(err).Msg("failed to send response")
			}
			return
		case chunk := <-out:
			if chunk.Err != "" {
				log.Error().Err(errors.New(chunk.Err)).Msg("failed to predict")
				if err := c.Send("Не удалось получить данные"); err != nil {
					log.Error().Err(err).Msg("failed to send response")
				}
				return
			}

			if chunk.DataType == "PREDICTION" && !chunk.Chunk {
				name := fmt.Sprintf("%d.json", time.Now().Unix())
				file, err := os.Create(name)
				if err != nil {
					log.Error().Err(err).Msg("failed to create file")
					return
				}
				mp := make(map[string]map[string]interface{})
				dataBytes, err := json.Marshal(chunk.Data)
				if err != nil {
					_ = file.Close()
					log.Error().Err(err).Msg("failed to marshal data")
					return
				}
				if err := json.Unmarshal(dataBytes, &mp); err != nil {
					_ = file.Close()
					log.Error().Err(err).Msg("failed to unmarshal data")
					return
				}
				isRegular, ok := mp["data"]["is_regular"].(bool)
				if !ok {
					_ = file.Close()
					log.Error().Msg("failed to get is_regular")
					return
				}
				if !isRegular {
					if err := c.Send("Это нерегулярная закупка. По ней не получится построить предсказание"); err != nil {
						log.Error().Err(err).Msg("failed to send response")
					}
					return
				}
				forecast, ok := mp["data"]["forecast"].([]any)
				if !ok {
					_ = file.Close()
					log.Error().Msg("failed to get forecast")
					return
				}
				if len(forecast) == 0 {
					if err := c.Send("Вы запросили период для предсказания, за который нет закупок"); err != nil {
						log.Error().Err(err).Msg("failed to send response")
					}
					return
				}
				outputJson, err := json.Marshal(mp["data"]["output_json"])
				if err != nil {
					_ = file.Close()
					log.Error().Err(err).Msg("failed to marshal output_json")
					return
				}
				_, err = file.Write(outputJson)
				if err != nil {
					_ = file.Close()
					log.Error().Err(err).Msg("failed to write data")
					return
				}
				_ = file.Close()

				if err := c.Send(&tele.Document{
					File:     tele.FromDisk(name),
					FileName: "prediction.json",
				}); err != nil {
					log.Error().Err(err).Msg("failed to send file")
				}
				continue
			}
			if chunk.Chunk {
				buff.WriteString(chunk.Data.(chunkMessage).Info)
				if buff.Len() > 350 {
					if err := c.Send(buff.String()); err != nil {
						log.Error().Err(err).Msg("failed to send response")
						return
					}
					buff.Reset()
				}
			}
			if chunk.Finish {
				return
			}
		}
	}
}

type chunkMessage struct {
	Info string
}
