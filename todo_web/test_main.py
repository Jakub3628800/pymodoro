from fastapi.testclient import TestClient
from main import app

client = TestClient(app)

def test_read_main():
    response = client.get("/")
    assert response.status_code == 401  # Unauthorized without credentials

def test_read_main_with_auth():
    response = client.get("/", auth=("user", "pass"))
    assert response.status_code == 200
    assert "Todo Lists" in response.text

def test_read_todo():
    response = client.get("/todo/test.json", auth=("user", "pass"))
    assert response.status_code == 200
    assert "test.json" in response.text

def test_update_todo():
    response = client.post("/todo/test.json", auth=("user", "pass"), data={"item_index": 0, "done": True})
    assert response.status_code == 200
    assert response.json() == {"success": True}

def test_read_main():
    response = client.get("/")
    assert response.status_code == 401  # Unauthorized without credentials

def test_read_main_with_auth():
    response = client.get("/", auth=("user", "pass"))
    assert response.status_code == 200
    assert "Todo Lists" in response.text

def test_read_todo():
    response = client.get("/todo/test.json", auth=("user", "pass"))
    assert response.status_code == 200
    assert "test.json" in response.text

def test_update_todo():
    response = client.post("/todo/test.json", auth=("user", "pass"), data={"item_index": 0, "done": True})
    assert response.status_code == 200
    assert response.json() == {"success": True}