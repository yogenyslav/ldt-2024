import { useContext } from 'react';
import { rootStoreContext } from '../stores';

export const useStores = () => useContext(rootStoreContext);
