import { createContext } from '@lit/context';

export const indexersContext = createContext<string[]>('indexers');
export const indexIdContext = createContext<string>('index-id');
