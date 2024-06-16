'use client';

import * as React from 'react';
import { Check, ChevronsUpDown } from 'lucide-react';

import { cn } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import {
    Command,
    CommandEmpty,
    CommandGroup,
    CommandInput,
    CommandItem,
    CommandList,
} from '@/components/ui/command';
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover';
import { Product } from '@/api/models/organizations';

type Props = {
    products: Product[];
};

export function ProductsComboBox({ products }: Props) {
    const [open, setOpen] = React.useState(false);
    const [value, setValue] = React.useState('');

    return (
        <div className='w-full'>
            <Popover open={open} onOpenChange={setOpen}>
                <PopoverTrigger asChild>
                    <Button
                        variant='outline'
                        role='combobox'
                        aria-expanded={open}
                        className='w-full justify-between'
                    >
                        {value
                            ? products.find((product) => product.name === value)?.name
                            : 'Поиск товара...'}
                        <ChevronsUpDown className='ml-2 h-4 w-4 shrink-0 opacity-50' />
                    </Button>
                </PopoverTrigger>
                <PopoverContent className='w-full max-w-3xl p-0'>
                    <Command>
                        <CommandInput placeholder='Поиск товара...' />
                        <CommandList>
                            <CommandEmpty>Ничего не найдено</CommandEmpty>
                            <CommandGroup>
                                {products.map((product) => (
                                    <CommandItem
                                        key={product.name}
                                        value={product.name}
                                        onSelect={(currentValue) => {
                                            setValue(currentValue === value ? '' : currentValue);
                                            setOpen(false);
                                        }}
                                    >
                                        <Check
                                            className={cn(
                                                'mr-2 h-4 w-4',
                                                value === product.name ? 'opacity-100' : 'opacity-0'
                                            )}
                                        />
                                        <span>
                                            {product.segment} {product.name}{' '}
                                            {product.regular ? 'Регулярная' : 'Нерегулярная'}
                                        </span>
                                    </CommandItem>
                                ))}
                            </CommandGroup>
                        </CommandList>
                    </Command>
                </PopoverContent>
            </Popover>
        </div>
    );
}
