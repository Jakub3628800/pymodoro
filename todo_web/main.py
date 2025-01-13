from fastapi import FastAPI, Request, Form, Depends, HTTPException, status
from fastapi.templating import Jinja2Templates
from fastapi.staticfiles import StaticFiles
from fastapi.security import HTTPBasic, HTTPBasicCredentials
from pydantic import BaseModel
import os
import json
from typing import List, Dict

app = FastAPI()
security = HTTPBasic()
templates = Jinja2Templates(directory="templates")
app.mount("/static", StaticFiles(directory="static"), name="static")

# Predefined credentials (in a real application, use more secure methods)
USERNAME = "user"
PASSWORD = "pass"

class TodoItem(BaseModel):
    task: str
    done: bool

def get_todo_files() -> List[str]:
    todo_dir = "/workspace/td"
    return [f for f in os.listdir(todo_dir) if f.endswith(".json")]

def read_todo_file(filename: str) -> List[TodoItem]:
    with open(os.path.join("/workspace/td", filename), "r") as f:
        data = json.load(f)
    return [TodoItem(**item) for item in data]

def write_todo_file(filename: str, items: List[TodoItem]):
    with open(os.path.join("/workspace/td", filename), "w") as f:
        json.dump([item.dict() for item in items], f, indent=2)

def authenticate(credentials: HTTPBasicCredentials = Depends(security)):
    if credentials.username != USERNAME or credentials.password != PASSWORD:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid credentials",
            headers={"WWW-Authenticate": "Basic"},
        )
    return credentials.username

@app.get("/")
async def root(request: Request, username: str = Depends(authenticate)):
    todo_files = get_todo_files()
    return templates.TemplateResponse("index.html", {"request": request, "todo_files": todo_files})

@app.get("/todo/{filename}")
async def get_todo(request: Request, filename: str, username: str = Depends(authenticate)):
    items = read_todo_file(filename)
    return templates.TemplateResponse("todo.html", {"request": request, "filename": filename, "items": items})

@app.post("/todo/{filename}")
async def update_todo(filename: str, item_index: int = Form(...), done: bool = Form(...), username: str = Depends(authenticate)):
    items = read_todo_file(filename)
    items[item_index].done = done
    write_todo_file(filename, items)
    return {"success": True}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)