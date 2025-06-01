import {buildSubgraphSchema} from '@apollo/subgraph';
import {ApolloServer} from '@apollo/server';
import {startStandaloneServer} from '@apollo/server/standalone';
import {dashboardServiceSchema} from './dashboard-service.graphql.js';
import {DashboardService} from './dashboard-service.js';

export async function run() {
    // ports are injects by Bit to `process.env.PORT`
    const port = parseInt(process.env.PORT, 10) || 3000;
    const dashboardService = DashboardService.from();

    const {typeDefs, resolvers} = dashboardServiceSchema(dashboardService);
    // The ApolloServer constructor requires two parameters: your schema
    // definition and your set of resolvers.
    const server = new ApolloServer({
        schema: buildSubgraphSchema({typeDefs, resolvers})
    });

    // Passing an ApolloServer instance to the `startStandaloneServer` function:
    //  1. creates an Express app
    //  2. installs your ApolloServer instance as middleware
    //  3. prepares your app to handle incoming requests
    const {url} = await startStandaloneServer(server, {
        listen: {
            port
        }
    });

    // eslint-disable-next-line no-console
    console.log(`🚀  Server ready at: ${url}`);

    return {
        port,
        stop: async () => {
            await server.stop();
        }
    }
}
