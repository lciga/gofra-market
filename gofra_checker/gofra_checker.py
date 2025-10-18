#!/usr/bin/env python3

import sys
import hashlib
import requests
import traceback

# Статусы чекера
def service_up():
    """Сервис работает корректно - флаг успешно размещен и получен"""
    print("[service is worked] - 101")
    exit(101)

def service_corrupt():
    """Сервис доступен, но работает некорректно - не удалось разместить/получить флаг"""
    print("[service is corrupt] - 102")
    exit(102)

def service_mumble():
    """Сервис не ответил вовремя или отвалился по таймауту"""
    print("[service is mumble] - 103")
    exit(103)

def service_down():
    """Сервис недоступен - порт закрыт или сервис мертв"""
    print("[service is down] - 104")
    exit(104)


def debug(err):
    """Для отладки можно раскомментировать"""
    pass
    # if isinstance(err, str):
    #     err = Exception(err)
    # traceback.print_exc()
    # raise err


class GobraChecker:
    def __init__(self, host: str, port: int = 8080):
        self.host = host
        self.port = port
        self.base_url = f"http://{host}:{port}/api"
        self.timeout = 5

    def generate_credentials(self, flag: str) -> tuple:
        """Генерация логина и пароля из флага как seed"""
        seed = hashlib.sha256(flag.encode()).hexdigest()
        username = f"user_{seed[:16]}"
        password = f"pass_{seed[16:32]}"
        return username, password

    def check_service_availability(self):
        """Проверка доступности сервиса"""
        try:
            resp = requests.get(f"http://{self.host}:{self.port}/", timeout=self.timeout)
            if resp.status_code >= 500:
                service_down()
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_down()

    def register_user(self, username: str, password: str):
        """Регистрация пользователя"""
        try:
            session = requests.Session()
            data = {
                "login": username,
                "password": password
            }
            resp = session.post(
                f"{self.base_url}/register",
                json=data,
                timeout=self.timeout
            )
            
            if resp.status_code == 201:
                return session
            elif resp.status_code >= 500:
                service_corrupt()
            else:
                # Если пользователь уже существует, пробуем залогиниться
                return self.login_user(username, password)
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def login_user(self, username: str, password: str):
        """Авторизация пользователя"""
        try:
            session = requests.Session()
            data = {
                "login": username,
                "password": password
            }
            resp = session.post(
                f"{self.base_url}/login",
                json=data,
                timeout=self.timeout
            )
            
            if resp.status_code == 200:
                return session
            elif resp.status_code == 401:
                debug(f"Login failed: 401 Unauthorized")
                service_corrupt()
            elif resp.status_code >= 500:
                service_corrupt()
            else:
                debug(f"Login failed: {resp.status_code}")
                service_corrupt()
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def get_user_info(self, session):
        """Получение информации о текущем пользователе"""
        try:
            resp = session.get(
                f"{self.base_url}/me",
                timeout=self.timeout
            )
            
            if resp.status_code == 200:
                return resp.json()
            else:
                debug(f"Get user info failed: {resp.status_code}")
                service_corrupt()
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def create_listing(self, session, flag: str) -> str:
        """Создание листинга с гофером и флагом в описании"""
        try:
            seed = hashlib.md5(flag.encode()).hexdigest()
            gofer_name = f"Gofer_{seed[:8]}"
            gofer_rarity = (int(seed[:2], 16) % 3) + 1  # Редкость от 1 до 3
            price = 10000  # Фиксированная цена
            
            data = {
                "gofer_name": gofer_name,
                "gofer_rarity": gofer_rarity,
                "price": price,
                "description": flag  # Флаг в описании
            }
            
            resp = session.post(
                f"{self.base_url}/listings",
                json=data,
                timeout=self.timeout
            )
            
            if resp.status_code == 201:
                result = resp.json()
                listing_id = result.get("id")
                if not listing_id:
                    debug("No listing id in response")
                    service_corrupt()
                return listing_id
            elif resp.status_code >= 500:
                service_corrupt()
            else:
                debug(f"Create listing failed: {resp.status_code}")
                service_corrupt()
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def get_my_listings(self, session) -> list:
        """Получение списка своих листингов"""
        try:
            resp = session.get(
                f"{self.base_url}/my-listings",
                timeout=self.timeout
            )
            
            if resp.status_code == 200:
                data = resp.json()
                return data.get("listings", [])
            else:
                debug(f"Get my listings failed: {resp.status_code}")
                service_corrupt()
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def get_my_gofers(self, session) -> list:
        """Получение списка своих гоферов"""
        try:
            resp = session.get(
                f"{self.base_url}/my-gofers",
                timeout=self.timeout
            )
            
            if resp.status_code == 200:
                data = resp.json()
                return data.get("gofers", [])
            else:
                debug(f"Get my gofers failed: {resp.status_code}")
                service_corrupt()
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def get_listing_by_id(self, session, listing_id: str):
        """Получение листинга по ID"""
        try:
            resp = session.get(
                f"{self.base_url}/listings/{listing_id}",
                timeout=self.timeout
            )
            
            if resp.status_code == 200:
                return resp.json()
            elif resp.status_code == 404:
                return None
            else:
                debug(f"Get listing failed: {resp.status_code}")
                service_corrupt()
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def search_market(self, session) -> list:
        """Проверка работы поиска на маркетплейсе"""
        try:
            resp = session.get(
                f"{self.base_url}/market",
                timeout=self.timeout
            )
            
            if resp.status_code == 200:
                data = resp.json()
                return data.get("items", [])
            else:
                debug(f"Market search failed: {resp.status_code}")
                service_corrupt()
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def buy_listing(self, session, listing_id: str):
        """Покупка листинга"""
        try:
            data = {
                "listing_id": listing_id
            }
            resp = session.post(
                f"{self.base_url}/buy",
                json=data,
                timeout=self.timeout
            )
            
            if resp.status_code == 204:
                return True
            elif resp.status_code >= 500:
                service_corrupt()
            else:
                # Может быть ошибка из-за недостатка средств или другая логическая ошибка
                debug(f"Buy listing failed: {resp.status_code}")
                return False
                
        except requests.exceptions.Timeout:
            service_mumble()
        except requests.exceptions.ConnectionError:
            service_down()
        except Exception as e:
            debug(e)
            service_corrupt()

    def put_flag(self, flag_id: str, flag: str):
        """Размещение флага в сервисе"""
        # 1. Проверяем доступность сервиса
        self.check_service_availability()
        
        # 2. Генерируем креденшиалы из флага
        username, password = self.generate_credentials(flag)
        
        # 3. Регистрируем пользователя
        session = self.register_user(username, password)
        if not session:
            service_corrupt()
        
        # 4. Проверяем, что можем получить информацию о пользователе
        user_info = self.get_user_info(session)
        if not user_info or "user_id" not in user_info:
            debug("No user_id in user info")
            service_corrupt()
        
        # 5. Создаем листинг с флагом в описании
        listing_id = self.create_listing(session, flag)
        if not listing_id:
            debug("No listing_id returned")
            service_corrupt()
        
        # 6. Проверяем, что листинг появился в списке своих листингов
        my_listings = self.get_my_listings(session)
        found = False
        for listing in my_listings:
            if listing.get("id") == listing_id:
                if listing.get("description") != flag:
                    debug(f"Flag mismatch in listing: expected '{flag}', got '{listing.get('description')}'")
                    service_corrupt()
                if listing.get("is_sold") == True:
                    debug("Listing is already sold after creation")
                    service_corrupt()
                found = True
                break
        
        if not found:
            debug(f"Listing {listing_id} not found in my listings")
            service_corrupt()
        
        # 7. Проверяем работу маркета (дополнительная проверка)
        market_items = self.search_market(session)
        # Маркет должен работать, даже если наш листинг там не виден
        
        # 8. Проверяем работу получения листинга по ID
        listing_data = self.get_listing_by_id(session, listing_id)
        if not listing_data:
            debug("Cannot get listing by id")
            service_corrupt()
        if listing_data.get("description") != flag:
            debug(f"Flag mismatch when getting by id: expected '{flag}', got '{listing_data.get('description')}'")
            service_corrupt()

    def check_flag(self, flag_id: str, flag: str):
        """Проверка наличия флага в сервисе"""
        # 1. Генерируем креденшиалы из флага (как seed)
        username, password = self.generate_credentials(flag)
        
        # 2. Проверяем доступность сервиса
        self.check_service_availability()
        
        # 3. Логинимся под пользователем
        session = self.login_user(username, password)
        if not session:
            service_corrupt()
        
        # 4. Получаем информацию о пользователе
        user_info = self.get_user_info(session)
        if not user_info or "user_id" not in user_info:
            debug("No user_id in user info during check")
            service_corrupt()
        
        user_id = user_info.get("user_id")
        
        # 5. Получаем список своих листингов
        my_listings = self.get_my_listings(session)
        
        # 6. Ищем листинг с нашим флагом
        found_flag = False
        flag_listing = None
        for lst in my_listings:
            if lst.get("description") == flag:
                # Проверяем, что листинг принадлежит нашему пользователю (seller_id)
                if lst.get("seller_id") != user_id:
                    debug(f"Listing seller_id mismatch: expected {user_id}, got {lst.get('seller_id')}")
                    service_corrupt()
                found_flag = True
                flag_listing = lst
                break
        
        if not found_flag:
            debug(f"Flag not found in listings for user {username}")
            service_corrupt()
        
        # 7. КРИТИЧНО: Проверяем, что листинг не был продан
        # Если листинг продан (is_sold == True), значит кто-то купил гофера с флагом
        # и флаг скомпрометирован -> сервис corrupt
        if flag_listing.get("is_sold") == True:
            debug(f"Flag compromised: listing {flag_listing.get('id')} was sold to {flag_listing.get('buyer_id')}")
            service_corrupt()
        
        # 8. Дополнительная проверка: получаем гоферов пользователя
        my_gofers = self.get_my_gofers(session)
        # У пользователя должен быть хотя бы один гофер (тот, что в листинге)
        if len(my_gofers) == 0:
            debug("User has no gofers")
            service_corrupt()
        
        # 9. Проверяем, что можем получить листинг по ID
        listing_id = flag_listing.get("id")
        listing_data = self.get_listing_by_id(session, listing_id)
        if not listing_data:
            debug(f"Cannot get listing {listing_id} by id during check")
            service_corrupt()
        
        # Проверяем, что описание всё ещё содержит флаг (т.к. мы владелец)
        if listing_data.get("description") != flag:
            debug(f"Flag mismatch in listing by id: expected '{flag}', got '{listing_data.get('description')}'")
            service_corrupt()
        
        # Еще раз проверяем is_sold
        if listing_data.get("is_sold") == True:
            debug(f"Flag compromised (from get by id): listing was sold")
            service_corrupt()


def main():
    if len(sys.argv) != 5:
        print("\nUsage:\n\t" + sys.argv[0] + " <host> (put|check) <flag_id> <flag>\n")
        print("Example:\n\t" + sys.argv[0] + ' "127.0.0.1" put "flag_id_123" "flag{test_flag_12345}"\n')
        exit(0)

    host = sys.argv[1]
    command = sys.argv[2]
    flag_id = sys.argv[3]
    flag = sys.argv[4]
    
    # Порт сервиса
    port = 8080
    
    checker = GobraChecker(host, port)
    
    try:
        if command == "put":
            # Размещаем флаг и сразу проверяем его
            checker.put_flag(flag_id, flag)
            checker.check_flag(flag_id, flag)
            service_up()
            
        elif command == "check":
            # Только проверяем наличие флага
            checker.check_flag(flag_id, flag)
            service_up()
            
        else:
            print(f"Unknown command: {command}")
            print("Available commands: put, check")
            exit(1)
            
    except KeyboardInterrupt:
        print("\nInterrupted by user")
        exit(1)
    except Exception as e:
        debug(f"Unexpected error: {e}")
        service_corrupt()


if __name__ == "__main__":
    main()
