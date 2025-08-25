import { createContext } from '@lit/context';

export const indexersContext = createContext<string[]>('indexers');

export const indexerIdContext = createContext<string>('indexer-id');

export interface Category {
  id: string;
  name: string;
  subCategories: Category[];
}

export interface IndexerDetails {
  categories: (indexer: string) => Promise<Category[]>;
}

export const indexerDetailsContext = createContext<IndexerDetails>('indexer-details');
