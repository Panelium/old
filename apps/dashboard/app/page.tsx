"use client";

import React from "react";
import {Button} from "@/components/ui/button.tsx";
import {TestService} from "@proto/common_pb.ts";
import {createClient} from "@connectrpc/connect";
import {createConnectTransport} from "@connectrpc/connect-web";

export default function Home() {
    const client = createClient(
        TestService,
        createConnectTransport({
            baseUrl: "http://localhost:8080",
        }),
    );

    async function handleTestMethod() {
        const response = await client.testMethod({
            text: "Hello from the dashboard app!",
            number: 42,
            boolean: true,
            array: ["item1", "item2", "item3"],
        });
        console.log("Response from testMethod:", response);
    }

    return (
        <div>
            <Button
                onClick={handleTestMethod}
                className="bg-blue-500 text-white hover:bg-blue-600 active:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
            >
                Test
            </Button>
        </div>
    );
}
