import { TodoItem } from '../types/todo';
import TodoCard from './TodoCard';

interface TodoListProps {
  todos: TodoItem[];
  onToggleComplete: (id: number, completed: boolean) => void;
  onDeleteTodo: (id: number) => void;
}

export default function TodoList({ todos, onToggleComplete, onDeleteTodo }: TodoListProps) {
  if (todos.length === 0) {
    return (
      <div className="bg-blue-50 p-8 rounded-lg text-center">
        <p className="text-blue-600 font-medium">タスクがありません。新しいタスクを追加してください。</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {todos.map((todo) => (
        <TodoCard
          key={todo.id}
          todo={todo}
          onToggleComplete={onToggleComplete}
          onDeleteTodo={onDeleteTodo}
        />
      ))}
    </div>
  );
} 