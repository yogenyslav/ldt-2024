import { useEffect, useState } from 'react';
import { ProductsComboBox } from '../ProductsComboBox';
import { Product } from '@/api/models/organizations';
import OrganizationsApiService from '@/api/OrganizationsApiService';

type Props = {
    organizationId: number;
};

const GoodsList = ({ organizationId }: Props) => {
    const [products, setProducts] = useState<Product[]>([]);

    useEffect(() => {
        OrganizationsApiService.getProducts(organizationId).then((products) => {
            setProducts(products.codes);
        });
    }, [organizationId]);

    return (
        <>
            <div className='mt-7'>
                <div className='flex flex-col'>
                    <h1 className='font-semibold text-lg md:text-2xl'>
                        Просмотр загруженных данных
                    </h1>

                    <p className='mt-2'>
                        Просмотр загруженных данных и поиск регулярных/нерегулярных закупок{' '}
                    </p>

                    <div className='mt-4'>
                        {products.length ? (
                            <ProductsComboBox products={products} />
                        ) : (
                            <p>Нет данных</p>
                        )}
                    </div>
                </div>
            </div>
        </>
    );
};

export default GoodsList;
