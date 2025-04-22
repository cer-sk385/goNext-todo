'use client';

import { useState, useEffect } from 'react';
import TodoForm from '../components/TodoForm';
import TodoList from '../components/TodoList';
import { TodoItem } from '../types/todo';

export default function Home() {
  const [todos, setTodos] = useState<TodoItem[]>([]);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  const fetchTodos = async () => {
    setIsLoading(true);
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/todos`);
      if (!response.ok) {
        throw new Error('データの取得に失敗しました');
      }
      const data = await response.json();
      setTodos(data);
      setError(null);
    } catch (err) {
      setError('Todoリストの読み込みに失敗しました');
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchTodos();
  }, []);

  const handleAddTodo = async (newTodo: Omit<TodoItem, 'id' | 'created_at' | 'updated_at'>) => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/todos`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newTodo),
      });

      if (!response.ok) {
        throw new Error('Todoの追加に失敗しました');
      }

      fetchTodos(); // リストを再取得して更新
    } catch (err) {
      setError('Todoの追加に失敗しました');
      console.error(err);
    }
  };

  const handleToggleComplete = async (id: number, completed: boolean) => {
    try {
      const todoToUpdate = todos.find(todo => todo.id === id);
      if (!todoToUpdate) return;

      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/todos/${id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ...todoToUpdate,
          completed: !completed
        }),
      });

      if (!response.ok) {
        throw new Error('Todoの更新に失敗しました');
      }

      fetchTodos(); // リストを再取得して更新
    } catch (err) {
      setError('Todoの更新に失敗しました');
      console.error(err);
    }
  };

  const handleDeleteTodo = async (id: number) => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/todos/${id}`, {
        method: 'DELETE',
      });

      if (!response.ok) {
        throw new Error('Todoの削除に失敗しました');
      }

      fetchTodos(); // リストを再取得して更新
    } catch (err) {
      setError('Todoの削除に失敗しました');
      console.error(err);
    }
  };

  return (
    <main className="min-h-screen bg-gradient-to-br from-purple-50 to-blue-50 p-4 md:p-8">
      <div className="max-w-3xl mx-auto">
        <div className="bg-white rounded-2xl shadow-xl overflow-hidden">
          <div className="bg-gradient-to-r from-indigo-500 to-purple-600 px-6 py-8">
            <h1 className="text-3xl font-bold text-white">TODOリスト</h1>
            <p className="text-indigo-100 mt-2">あなたのタスクを整理しましょう</p>
          </div>
          
          <div className="p-6">
            <TodoForm onAddTodo={handleAddTodo} />
            
            {error && (
              <div className="bg-red-50 text-red-600 p-4 rounded-lg my-4">
                {error}
              </div>
            )}
            
            {isLoading ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-500"></div>
              </div>
            ) : (
              <TodoList 
                todos={todos} 
                onToggleComplete={handleToggleComplete}
                onDeleteTodo={handleDeleteTodo}
              />
            )}
          </div>
        </div>
      </div>
    </main>
  );
} 