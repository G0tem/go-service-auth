import asyncio
import aiohttp
import time
from .utils import (
    plot_response_times,
    save_results_to_file,
    output_analitics,
    generate_user_data,
    plot_performance_comparison,
    plot_errors_over_time
)
from .settings import BASE_URL, REGISTER_ENDPOINT, LOGIN_ENDPOINT, ENDPOINTS, NUM_USERS


async def register_user(session, user_data):
    """
    Регистрация пользователя.
    """
    headers = {"accept": "application/json", "Content-Type": "application/json"}
    try:
        start_time = time.time()
        async with session.post(
            f"{BASE_URL}{REGISTER_ENDPOINT}", json=user_data, headers=headers
        ) as response:
            elapsed_time = time.time() - start_time
            if response.status != 201:
                return {
                    "status": "error",
                    "time": elapsed_time,
                    "code": response.status,
                    "endpoint": REGISTER_ENDPOINT,
                    "message": await response.text(),
                }
            data = await response.json()
            token = data.get("data", {}).get("token")
            return {
                "status": "success",
                "time": elapsed_time,
                "code": response.status,
                "endpoint": REGISTER_ENDPOINT,
                "token": token,
            }
    except Exception as e:
        elapsed_time = time.time() - start_time
        return {
            "status": "error",
            "time": elapsed_time,
            "code": 500,
            "endpoint": REGISTER_ENDPOINT,
            "message": str(e),
        }


async def login_user(session, user_data):
    """
    Авторизация пользователя.
    """
    headers = {"accept": "application/json", "Content-Type": "application/json"}
    login_data = {"identity": user_data["email"], "password": user_data["password"]}
    try:
        start_time = time.time()
        async with session.post(
            f"{BASE_URL}{LOGIN_ENDPOINT}", json=login_data, headers=headers
        ) as response:
            elapsed_time = time.time() - start_time
            if response.status != 200:
                return {
                    "status": "error",
                    "time": elapsed_time,
                    "code": response.status,
                    "endpoint": LOGIN_ENDPOINT,
                    "message": await response.text(),
                }
            data = await response.json()
            token = data.get("data", {}).get(
                "token"
            )
            return {
                "status": "success",
                "time": elapsed_time,
                "code": response.status,
                "endpoint": LOGIN_ENDPOINT,
                "token": token,
            }
    except Exception as e:
        elapsed_time = time.time() - start_time
        return {
            "status": "error",
            "time": elapsed_time,
            "code": 500,
            "endpoint": LOGIN_ENDPOINT,
            "message": str(e),
        }


async def make_request(session, endpoint, token):
    """
    Выполнение GET-запроса.
    """
    start_time = time.time()
    headers = {"accept": "application/json", "Authorization": f"Bearer {token}"}
    try:
        async with session.get(f"{BASE_URL}{endpoint}", headers=headers) as response:
            elapsed_time = time.time() - start_time
            if response.status != 200:
                return {
                    "status": "error",
                    "time": elapsed_time,
                    "code": response.status,
                    "endpoint": endpoint,
                }
            return {
                "status": "success",
                "time": elapsed_time,
                "code": response.status,
                "endpoint": endpoint,
            }
    except Exception as e:
        elapsed_time = time.time() - start_time
        return {
            "status": "error",
            "time": elapsed_time,
            "message": str(e),
            "endpoint": endpoint,
        }


async def user_workflow(user_data):
    """
    Полный цикл работы пользователя: регистрация, авторизация и запросы на эндпоинты.
    """
    async with aiohttp.ClientSession() as session:
        results = []

        # Регистрация пользователя
        registration_result = await register_user(session, user_data)
        results.append(registration_result)

        # Авторизация пользователя
        if registration_result["status"] == "success":
            login_result = await login_user(session, user_data)
            results.append(login_result)

            # endpoint из списка
            if login_result["status"] == "success":
                for endpoint in ENDPOINTS:
                    result = await make_request(
                        session, endpoint, login_result["token"]
                    )
                    results.append(result)

        return results


async def run_tests():
    """
    Основная функция для запуска тестов.
    """
    # Генерация данных для регистрации
    user_data_list = [generate_user_data(i) for i in range(NUM_USERS)]

    # Запуск полного цикла для каждого пользователя
    tasks = [user_workflow(user_data) for user_data in user_data_list]
    all_results = await asyncio.gather(*tasks)

    # Объединение результатов всех пользователей
    results = [result for user_results in all_results for result in user_results]

    # Формирование аналитики
    analytics = output_analitics(results)

    plot_response_times(results)
    plot_performance_comparison(results)

    save_results_to_file(results, analytics)

    plot_errors_over_time(results)

def test_performance():
    asyncio.run(run_tests())
