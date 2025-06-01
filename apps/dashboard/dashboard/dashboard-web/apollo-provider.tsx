import {ReactNode} from 'react';
import {ApolloClient, ApolloProvider, InMemoryCache} from '@apollo/client';

export type BeerShopApolloProviderProps = {
    children: ReactNode;
};

export function AcmeApolloProvider({children}: BeerShopApolloProviderProps) {
    const gatewayUrl = process.env?.BACKEND_URL;

    const client = new ApolloClient({
        uri: gatewayUrl,
        cache: new InMemoryCache(),
        credentials: 'include'
    });

    return (
        <ApolloProvider client={client}>
            {children}
        </ApolloProvider>
    )
}
