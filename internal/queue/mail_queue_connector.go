package queue

import (
	"encoding/json"
	"time"

	"github.com/G0tem/go-servise-auth/internal/config"
	"github.com/G0tem/go-servise-auth/internal/default_email_templates"
	"github.com/G0tem/go-servise-auth/internal/dto"
	"github.com/G0tem/go-servise-auth/internal/model"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

type MailQueueConnector struct {
	cfg *config.Config
}

func (qConn *MailQueueConnector) getChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	// No defer, channel will automatically closed on error

	err = ch.ExchangeDeclarePassive(
		qConn.cfg.RMQMailExchange, // name
		"fanout",                  // type
		true,                      // durable
		false,                     // auto-deleted
		false,                     // internal
		false,                     // no-wait
		nil,                       // arguments
	)
	if err != nil && qConn.cfg.RMQMailExchangeAutocreate {
		log.Warn().Msgf("Exchange %q doesn't exist, autocreate enabled, so attempting to create exchange", qConn.cfg.RMQMailExchange)

		autocreateExchange(conn, qConn.cfg.RMQMailExchange)

		ch, err = conn.Channel()
		if err != nil {
			return nil, err
		}
	}

	return ch, nil
}

func (qConn *MailQueueConnector) PushEmailConfirm(user *model.User, token string) error {
	emailTemplate := model.NewEmailTemplate(default_email_templates.EMAIL_CONFIRM)
	smtpMessage, err := emailTemplate.Apply(dto.MailHtml, user.Email, struct {
		Title string
		Ref   string
		Name  string
	}{
		Title: user.Username,
		Ref:   qConn.cfg.PublicEmailConfirmationUrl + "?" + token,
		Name:  user.Username,
	})
	if err != nil {
		return err
	}

	body, err := json.Marshal(smtpMessage)
	if err != nil {
		return err
	}

	conn, err := amqp.Dial(qConn.cfg.RMQConnUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := qConn.getChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		qConn.cfg.RMQMailExchange, // exchange
		"mail.*",                  // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			Timestamp:   time.Now(),
		})
	if err != nil {
		return err
	}

	return nil
}

func (qConn *MailQueueConnector) PushPasswordResetConfirm(user *model.User, token string) error {
	conn, err := amqp.Dial(qConn.cfg.RMQConnUrl)
	if err != nil {
		return err
	}
	defer conn.Close()

	emailTemplate := model.NewEmailTemplate(default_email_templates.PASSWORD_RESET_CONFIRM)
	smtpMessage, err := emailTemplate.Apply(dto.MailHtml, user.Email, struct {
		Title string
		Ref   string
		Name  string
	}{
		Title: user.Username,
		Ref:   qConn.cfg.PublicPasswordResetConfirmationUrl + token,
		Name:  user.Username,
	})
	if err != nil {
		return err
	}

	body, err := json.Marshal(smtpMessage)
	if err != nil {
		return err
	}

	ch, err := qConn.getChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()

	err = ch.Publish(
		qConn.cfg.RMQMailExchange, // exchange
		"mail.*",                  // routing key
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
			Timestamp:   time.Now(),
		})
	if err != nil {
		return err
	}

	return nil
}

func NewMailQueueConnector(cfg *config.Config) *MailQueueConnector {
	return &MailQueueConnector{
		cfg: cfg,
	}
}
