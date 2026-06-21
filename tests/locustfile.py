import random
import string
import time
from locust import HttpUser, task, between


class ApiUser(HttpUser):
    wait_time = between(1, 3)

    def on_start(self):
        self._new_identity()

    def _new_identity(self):
        suffix = "".join(random.choices(string.ascii_lowercase + string.digits, k=8))
        ts = str(int(time.time() * 1000))[-6:]
        self.username = f"loadtest_{suffix}_{ts}"
        self.password = "TestPass123!"
        self.email = f"{self.username}@test.com"
        self.full_name = f"Load Test {self.username}"
        self.token = None
        self.logged_in = False
        self._signup()

    def _signup(self):
        with self.client.post(
            "/api/signup",
            json={
                "username": self.username,
                "password": self.password,
                "email": self.email,
                "full_name": self.full_name,
            },
            catch_response=True,
        ) as resp:
            if resp.status_code == 200 and resp.json().get("message") == "user registered successfully":
                resp.success()
            else:
                resp.failure(f"Signup failed: {resp.status_code} {resp.text}")

    @task(3)
    def login(self):
        with self.client.post(
            "/api/login",
            json={"username": self.username, "password": self.password},
            catch_response=True,
        ) as resp:
            if resp.status_code == 200:
                data = resp.json()
                if data.get("token"):
                    self.token = data["token"]
                    self.logged_in = True
                    resp.success()
                else:
                    resp.failure("Login response missing token")
            else:
                resp.failure(f"Login failed: {resp.status_code} {resp.text}")

    @task(5)
    def profile(self):
        if not self.logged_in or not self.token:
            self.login()
            if not self.token:
                return

        with self.client.get(
            "/api/profile",
            headers={"Authorization": f"Bearer {self.token}"},
            catch_response=True,
        ) as resp:
            if resp.status_code == 200:
                data = resp.json().get("data", {})
                if data.get("username") == self.username:
                    resp.success()
                else:
                    resp.failure(f"Profile username mismatch: expected {self.username}, got {data.get('username')}")
            elif resp.status_code in (401, 400):
                self.logged_in = False
                self.token = None
                resp.failure(f"Profile unauthorized, will re-login: {resp.status_code}")
            else:
                resp.failure(f"Profile failed: {resp.status_code} {resp.text}")
