'use client';

import { useState } from 'react';
import { TodoItem } from '../types/todo';

interface TodoFormProps {
  onAddTodo: (todo: Omit<TodoItem, 'id' | 'created_at' | 'updated_at'>) => void;
}

export default function TodoForm({ onAddTodo }: TodoFormProps) {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [isExpanded, setIsExpanded] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!title.trim()) return;
    
    onAddTodo({
      title,
      description,
      completed: false
    });
    
    // フォームをリセット
    setTitle('');
    setDescription('');
    setIsExpanded(false);
  };

  return (
    <form onSubmit={handleSubmit} className="mb-8">
      <div className="flex flex-col gap-4">
        <div>
          <div className="flex items-center gap-2 bg-gray-50 rounded-lg border border-gray-200 p-2 focus-within:ring-2 focus-within:ring-indigo-500 focus-within:border-transparent">
            <input
              type="text"
              placeholder="新しいタスクを追加"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              onFocus={() => setIsExpanded(true)}
              className="flex-1 bg-transparent outline-none px-2 py-1"
            />
            <button
              type="submit"
              disabled={!title.trim()}
              className="bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              追加
            </button>
          </div>
        </div>
        
        {isExpanded && (
          <div className="transition-all duration-300">
            <textarea
              placeholder="詳細説明（任意）"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              className="w-full p-3 bg-gray-50 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent min-h-[100px]"
            />
          </div>
        )}
      </div>
    </form>
  );
} 