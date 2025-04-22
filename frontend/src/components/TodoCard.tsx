import { TodoItem } from '../types/todo';
import { formatDistanceToNow, parseISO } from 'date-fns';
import { ja } from 'date-fns/locale';

interface TodoCardProps {
  todo: TodoItem;
  onToggleComplete: (id: number, completed: boolean) => void;
  onDeleteTodo: (id: number) => void;
}

export default function TodoCard({ todo, onToggleComplete, onDeleteTodo }: TodoCardProps) {
  const formattedDate = formatDistanceToNow(parseISO(todo.updated_at), {
    addSuffix: true,
    locale: ja,
    includeSeconds: true
  });

  return (
    <div className={`border rounded-xl p-4 shadow-sm transition-all ${
      todo.completed ? 'bg-gray-50 border-gray-200' : 'bg-white border-gray-200 hover:border-indigo-300'
    }`}>
      <div className="flex items-start gap-3">
        <button
          onClick={() => onToggleComplete(todo.id, todo.completed)}
          className="mt-1 flex-shrink-0"
        >
          <div className={`w-5 h-5 rounded-full border flex items-center justify-center transition-colors ${
            todo.completed
              ? 'bg-green-500 border-green-500 text-white'
              : 'border-gray-300 hover:border-indigo-500'
          }`}>
            {todo.completed && (
              <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" className="w-3 h-3">
                <path fillRule="evenodd" d="M16.704 4.153a.75.75 0 0 1 .143 1.052l-8 10.5a.75.75 0 0 1-1.127.075l-4.5-4.5a.75.75 0 0 1 1.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 0 1 1.05-.143Z" clipRule="evenodd" />
              </svg>
            )}
          </div>
        </button>
        
        <div className="flex-1 min-w-0">
          <div className="flex justify-between items-start">
            <h3 className={`font-medium text-lg ${
              todo.completed ? 'text-gray-500 line-through' : 'text-gray-900'
            }`}>
              {todo.title}
            </h3>
            <button
              onClick={() => onDeleteTodo(todo.id)}
              className="text-gray-400 hover:text-red-500 transition-colors p-1"
            >
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="w-5 h-5">
                <path strokeLinecap="round" strokeLinejoin="round" d="m14.74 9-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 0 1-2.244 2.077H8.084a2.25 2.25 0 0 1-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 0 0-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 0 1 3.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 0 0-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 0 0-7.5 0" />
              </svg>
            </button>
          </div>
          
          {todo.description && (
            <p className={`mt-1 text-sm ${
              todo.completed ? 'text-gray-400' : 'text-gray-600'
            }`}>
              {todo.description}
            </p>
          )}
          
          <p className="text-xs text-gray-400 mt-2">
            {formattedDate}
          </p>
        </div>
      </div>
    </div>
  );
} 