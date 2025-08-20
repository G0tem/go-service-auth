import matplotlib.pyplot as plt
import random
import string
from .settings import (
    REGISTER_ENDPOINT,
    LOGIN_ENDPOINT,
    ENDPOINTS,
)

def plot_performance_comparison(results):
    """
    Построение графика сравнения производительности этапов.
    """
    # Группировка результатов по этапам
    registration_times = [r["time"] for r in results if r["endpoint"] == REGISTER_ENDPOINT]
    login_times = [r["time"] for r in results if r["endpoint"] == LOGIN_ENDPOINT]
    endpoint_times = [r["time"] for r in results if r["endpoint"] in ENDPOINTS]

    # Подготовка данных
    stages = ["Registration", "Login", "Endpoints"]
    avg_times = [
        sum(registration_times) / len(registration_times) if registration_times else 0,
        sum(login_times) / len(login_times) if login_times else 0,
        sum(endpoint_times) / len(endpoint_times) if endpoint_times else 0,
    ]

    # Построение графика
    plt.figure(figsize=(8, 6))
    plt.bar(stages, avg_times, color=["blue", "orange", "green"])
    plt.title("Среднее время ответа")
    plt.xlabel("Этапы")
    plt.ylabel("Среднее время отклика (секунды)")
    plt.grid(axis='y', linestyle='--', alpha=0.7)

    plt.savefig("performance_comparison.png")
    plt.close()

def plot_response_times(results):
    """
    Построение графиков распределения времени ответа по эндпоинтам.
    """
    # Группировка результатов по эндпоинтам
    endpoint_times = {}
    for result in results:
        endpoint = result.get("endpoint", "Unknown")
        if endpoint not in endpoint_times:
            endpoint_times[endpoint] = []
        endpoint_times[endpoint].append(result["time"])

    # Построение графиков
    plt.figure(figsize=(12, 6))
    for i, (endpoint, times) in enumerate(endpoint_times.items(), start=1):
        plt.subplot(1, len(endpoint_times), i)
        plt.hist(times, bins=20, color='blue', edgecolor='black', alpha=0.7)
        plt.title(f"Endpoint: {endpoint}\nRequests: {len(times)}")
        plt.xlabel("Время отклика (секунды)")
        plt.ylabel("Количество запросов")
        plt.grid(axis='y', linestyle='--', alpha=0.7)

    plt.tight_layout()
    plt.savefig("response_time_distribution.png")
    plt.close()

def save_results_to_file(results, analytics, filename="test_results.txt"):
    """
    Сохранение результатов тестов в текстовый файл.
    Сначала записывается общая статистика, затем подробные результаты.
    """
    with open(filename, "w", encoding="utf-8") as file:
        # Запись общей аналитики
        file.write("Общая статистика:\n")
        file.write("=" * 50 + "\n")
        for line in analytics.split("\n"):
            file.write(line + "\n")
        file.write("=" * 50 + "\n\n")

        # Запись подробных результатов
        file.write("Подробные результаты тестов:\n")
        file.write("=" * 50 + "\n")
        for i, result in enumerate(results, start=1):
            status = result["status"]
            time_taken = result["time"]
            code = result.get("code", "N/A")
            endpoint = result.get("endpoint", "Unknown")  # Получаем эндпоинт из результата
            error_message = result.get("message", "")
            
            file.write(f"Запрос #{i}:\n")
            file.write(f"  Статус: {status}\n")
            file.write(f"  Время выполнения: {time_taken:.4f} секунд\n")
            file.write(f"  Код ответа: {code}\n")
            file.write(f"  Эндпоинт: {endpoint}\n")
            if status == "error":
                file.write(f"  Сообщение об ошибке: {error_message}\n")
            file.write("-" * 50 + "\n")

def output_analitics(results):
    registration_results = [r for r in results if r["endpoint"] == REGISTER_ENDPOINT]
    login_results = [r for r in results if r["endpoint"] == LOGIN_ENDPOINT]
    endpoint_results = [r for r in results if r["endpoint"] in ENDPOINTS]

    analytics = []

    analytics.append("\nАналитика регистрации:")
    analytics.append(print_metrics(registration_results))

    analytics.append("\nАналитика авторизации:")
    analytics.append(print_metrics(login_results))

    analytics.append("\nАналитика эндпоинтов:")
    analytics.append(print_metrics(endpoint_results))

    return "\n".join(analytics)

def print_metrics(results):
    total_time = sum(result["time"] for result in results)
    success_count = sum(1 for result in results if result["status"] == "success")
    error_count = len(results) - success_count
    avg_time = total_time / len(results) if results else 0
    max_time = max(result["time"] for result in results) if results else 0

    metrics = [
        f"Всего запросов: {len(results)}",
        f"Успешно: {success_count / len(results) * 100:.2f}%" if results else "Нет данных",
        f"Среднее время ответа: {avg_time:.4f} seconds",
        f"Максимальное время ответа: {max_time:.4f} seconds",
        f"Ошибки: {error_count}",
    ]
    return "\n".join(metrics)

def generate_random_string(length=8):
    """
    Функция для генерации случайных строк
    """
    return ''.join(random.choices(string.ascii_lowercase + string.digits, k=length))

def generate_user_data(index):
    """
    Функция для генерации уникальных данных для регистрации
    """
    username = f"user_{index}_{generate_random_string()}"
    email = f"user_{index}_{generate_random_string()}@example.com"
    password = "string"  # Одинаковый пароль для всех
    return {
        "confirmPassword": password,
        "email": email,
        "password": password,
        "username": username
    }

def plot_errors_over_time(results):
    """
    Построение графика ошибок по времени.
    """
    # Сортировка результатов по времени
    sorted_results = sorted(results, key=lambda x: x["time"])
    error_counts = []
    timestamps = []

    current_errors = 0
    for result in sorted_results:
        if result["status"] == "error":
            current_errors += 1
        error_counts.append(current_errors)
        timestamps.append(result["time"])

    # Построение графика
    plt.figure(figsize=(10, 6))
    plt.plot(timestamps, error_counts, color="red", label="Ошибки")
    plt.title("Ошибки по времени")
    plt.xlabel("Время (сукунды)")
    plt.ylabel("Кол-во ошибок")
    plt.grid(axis='y', linestyle='--', alpha=0.7)
    plt.legend()

    plt.savefig("errors_over_time.png")
    plt.close()