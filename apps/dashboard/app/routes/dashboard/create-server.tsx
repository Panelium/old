import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { Card } from "~/components/ui/card";
import { Client } from "@connectrpc/connect";
import { getClientClient } from "~/lib/api-clients";
import {
  AvailableBlueprint,
  AvailableLocation,
  AvailableNode,
  ClientService,
  NewServerRequest,
} from "proto-gen-ts/backend/Client_pb";
import { FormControl, FormField, FormItem, FormLabel, FormMessage } from "~/components/ui/form";
import { FormProvider, useForm } from "react-hook-form";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "~/components/ui/select";
import { Button } from "~/components/ui/button";
import { Input } from "~/components/ui/input";
import { Textarea } from "~/components/ui/textarea";

export default function ServerCreatePage() {
  const navigate = useNavigate();

  const [clientClient, setClientClient] = useState<Client<typeof ClientService>>();
  const [availableBlueprints, setAvailableBlueprints] = useState<AvailableBlueprint[]>([]);
  const [availableLocations, setAvailableLocations] = useState<AvailableLocation[]>([]);
  const [availableNodes, setAvailableNodes] = useState<AvailableNode[]>([]);

  const form = useForm({
    defaultValues: {
      name: "",
      description: "",
      bid: "",
      lid: "",
      nid: "",
    },
  });

  useEffect(() => {
    (async () => {
      const clientClient = await getClientClient();
      setClientClient(clientClient);

      const blueprints = await clientClient.getAvailableBlueprints({});
      setAvailableBlueprints(blueprints.blueprints);
      if (blueprints.blueprints.length > 0) {
        form.setValue("bid", blueprints.blueprints[0].bid);
      }
      const locations = await clientClient.getAvailableLocations({});
      setAvailableLocations(locations.locations);
      const nodes = await clientClient.getAvailableNodes({});
      setAvailableNodes(nodes.nodes);
    })();
  }, []);

  const refreshData = async () => {
    if (!clientClient) return;

    try {
      const blueprints = await clientClient.getAvailableBlueprints({});
      setAvailableBlueprints(blueprints.blueprints);
      const locations = await clientClient.getAvailableLocations({});
      setAvailableLocations(locations.locations);
      const nodes = await clientClient.getAvailableNodes({});
      setAvailableNodes(nodes.nodes);
    } catch (error) {
      console.error("Failed to refresh data:", error);
    }
  };

  const handleSubmit = form.handleSubmit(async (data) => {
    const payload: NewServerRequest = {
      $typeName: "backend.NewServerRequest",
      name: data.name,
      description: data.description,
      bid: data.bid,
    };
    if (data.lid && data.lid !== "__none__" && data.lid !== "__any__") payload.lid = data.lid;
    if (data.nid && data.nid !== "__none__" && data.nid !== "__any__") payload.nid = data.nid;

    const res = await clientClient?.newServer(payload);
    if (!res || !res.sid) {
      console.error("Failed to create server");
      return;
    }

    navigate(`/server/${res.sid}`, { replace: true });
  });

  return (
    <div className="flex justify-center items-center min-h-[60vh]">
      <Card className="w-full max-w-md p-8">
        <h2 className="text-2xl font-bold mb-4">Create New Server</h2>
        <FormProvider {...form}>
          <form onSubmit={handleSubmit} className="flex flex-col gap-4">
            <FormField
              control={form.control}
              name="name"
              rules={{ required: "Server name is required" }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Server Name</FormLabel>
                  <FormControl>
                    <Input type="text" className="input input-bordered" placeholder="Hypixel" {...field} />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="description"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Description</FormLabel>
                  <FormControl>
                    <Textarea
                      className="input input-bordered"
                      placeholder="A Minecraft server that hosts minigames."
                      {...field}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="bid"
              rules={{ required: "Blueprint is required" }}
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Blueprint</FormLabel>
                  <FormControl>
                    <Select value={field.value} onValueChange={field.onChange} defaultValue="">
                      <SelectTrigger className="input input-bordered w-full">
                        <SelectValue placeholder="Select a blueprint" />
                      </SelectTrigger>
                      <SelectContent>
                        {availableBlueprints.length === 0 ? (
                          <SelectItem value="__none__" disabled>
                            No blueprints available
                          </SelectItem>
                        ) : (
                          availableBlueprints.map((bp) => (
                            <SelectItem key={bp.bid} value={bp.bid}>
                              {bp.name}
                            </SelectItem>
                          ))
                        )}
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="lid"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Location</FormLabel>
                  <FormControl>
                    <Select value={field.value} onValueChange={field.onChange} defaultValue="">
                      <SelectTrigger className="input input-bordered w-full">
                        <SelectValue placeholder="Any" />
                      </SelectTrigger>
                      <SelectContent>
                        {availableLocations.length === 0 ? (
                          <SelectItem value="__none__" disabled>
                            No locations available
                          </SelectItem>
                        ) : (
                          <>
                            <SelectItem value="__any__">Any</SelectItem>
                            {availableLocations.map((loc) => (
                              <SelectItem key={loc.lid} value={loc.lid}>
                                {loc.name}
                              </SelectItem>
                            ))}
                          </>
                        )}
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <FormField
              control={form.control}
              name="nid"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Node</FormLabel>
                  <FormControl>
                    <Select value={field.value} onValueChange={field.onChange} defaultValue="">
                      <SelectTrigger className="input input-bordered w-full">
                        <SelectValue placeholder="Any" />
                      </SelectTrigger>
                      <SelectContent>
                        {availableNodes.length === 0 ? (
                          <SelectItem value="__none__" disabled>
                            No nodes available
                          </SelectItem>
                        ) : (
                          <>
                            <SelectItem value="__any__">Any</SelectItem>
                            {availableNodes
                              .filter(
                                (node) =>
                                  node.nid !== "__none__" &&
                                  node.nid !== "__any__" &&
                                  (form.getValues().bid && form.getValues().bid != ""
                                    ? form.getValues().bid === node.lid
                                    : true)
                              )
                              .map((node) => (
                                <SelectItem key={node.nid} value={node.nid}>
                                  {node.name}
                                </SelectItem>
                              ))}
                          </>
                        )}
                      </SelectContent>
                    </Select>
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit" className="btn btn-primary mt-2">
              Create Server
            </Button>
          </form>
        </FormProvider>
      </Card>
    </div>
  );
}
