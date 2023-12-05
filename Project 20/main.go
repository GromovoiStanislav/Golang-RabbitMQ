package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Message структура для сообщения
type Message struct {
	Text string `json:"text"`
}

type PackageMessages struct {
	Messages []string `json:"messages"`
}

type Message1C struct {
	Routing_key     string          `json:"routing_key"`
	PackageMessages PackageMessages `json:"package_messages"`
	Exchange        string          `json:"exchange"`
	Queue           string          `json:"queue"`
}

var amqpServerURL string 

func main() {
	// Загрузите переменные среды из файла .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("%s: %s", "Error loading .env file", err)
	}

	amqpServerURL = os.Getenv("AMQP_SERVER_URL")
	if amqpServerURL == "" {
		amqpServerURL  = "amqp://test:test@192.168.0.152:5672/"
	}

	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/get", handleGetMessage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleGetMessage(w http.ResponseWriter, r *http.Request) {

	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Только запросы POST", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса.", http.StatusInternalServerError)
		return
	}

	// Получаем сообщения из RabbitMQ
	var mesFrom1c Message1C

	if err := json.Unmarshal([]byte(string(body)), &mesFrom1c); err != nil {
		panic(err)
	}
	messages, err := getMessagesFromRabbitMQ(mesFrom1c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Преобразуем сообщения в JSON-массив
	jsonData, err := json.Marshal(messages)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-массив в теле ответа
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func getMessagesFromRabbitMQ(mes1C Message1C) (PackageMessages, error) {

	var packageMess PackageMessages

	// Подключаемся к RabbitMQ
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		return packageMess, fmt.Errorf("Ошибка подключения к RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Создаем канал
	ch, err := conn.Channel()
	if err != nil {
		return packageMess, fmt.Errorf("Ошибка открытия канала: %v", err)
	}
	defer ch.Close()
	fmt.Printf(mes1C.Queue)
	// Объявляем очередь
	queue, err := ch.QueueDeclare(
		mes1C.Queue, // Имя очереди
		true,        // Устойчивая ли очередь
		false,       // Автоматическое удаление при отключении последнего потребителя
		false,       // Эксклюзивная ли очередь
		false,       // Нужно ли ожидать подтверждения при передаче сообщения
		nil,         // Аргументы
	)
	if err != nil {
		return packageMess, fmt.Errorf("Ошибка объявления очереди: %v", err)
	}

	// Создаем канал для остановки
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	// Получаем сообщения из очереди
	msgs, err := ch.Consume(
		mes1C.Queue, // Имя очереди
		"",          // Идентификатор потребителя
		true,        // Автоматическое подтверждение
		false,       // Эксклюзивный ли потребитель
		false,       // Нужно ли ждать подтверждения при передаче сообщения
		false,       // Нужно ли подтверждать получение сообщений
		nil,         // Аргументы
	)
	if err != nil {
		fmt.Printf(string(err.Error()))
		return packageMess, fmt.Errorf("Ошибка регистрации получателя: %v", err)
	}
	fmt.Printf("4")

	// Собираем сообщения в массив
	msgCount := 0
	fmt.Println("Queue length:", queue.Messages)
	for {
		if msgCount == queue.Messages {
			return packageMess, nil
		}

		select {

		case d, ok := <-msgs:
			if !ok {
				// Канал закрыт
				return packageMess, nil
			}
			packageMess.Messages = append(packageMess.Messages, string(d.Body))
			msgCount++
		case <-stopChan:
			// Получен сигнал остановки
			return packageMess, nil
		}
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Только запросы POST.", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения тела запроса.", http.StatusInternalServerError)
		return
	}

	// Указываем адрес сервера RabbitMQ
	conn, err := amqp.Dial(amqpServerURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к RabbitMQ: %v", err)
	}
	defer conn.Close()

	start := time.Now()

	// Создаем канал
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала: %v", err)
	}
	defer ch.Close()

	var mesFrom1c Message1C

	if err := json.Unmarshal([]byte(string(body)), &mesFrom1c); err != nil {
		panic(err)
	}

	fmt.Println(len(mesFrom1c.PackageMessages.Messages))

	for _, str := range mesFrom1c.PackageMessages.Messages {
		err = ch.Publish(
			mesFrom1c.Exchange,    // exchange
			mesFrom1c.Routing_key, // routing key
			false,                 // mandatory
			false,                 // immediate
			amqp.Publishing{
				ContentType: "text/json",
				Body:        []byte(str),
			})
	}

	if err != nil {
		log.Fatalf("Ошибка публикации сообщения: %v", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Время выполнения: %s\n", elapsed.Seconds())
}
