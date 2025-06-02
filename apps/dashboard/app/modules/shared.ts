// This is a placeholder for a shared utility library for the frontend modules.
// For example, functions for API calls, event bus, or shared UI components specific to modules.

export const sayHelloFromSharedModuleLib = () => {
    console.log("Hello from the shared module library!");
};

// Interface for a generic module
export interface IPaneliumModule {
    name: string;
    version: string;
    // Routes provided by the module (e.g., for admin settings or user-facing pages)
    routes?: {
        user?: any[]; // Define a more specific type for route configuration
        admin?: any[];
    };
    // Components to be rendered in specific parts of the UI (e.g., dashboard widgets, settings tabs)
    components?: {
        userDashboardWidgets?: React.ComponentType[];
        adminSettingsTabs?: React.ComponentType[];
        // ... other extension points
    };
    // Functions to register with the core system (e.g., event listeners, service providers)
    register?: (coreAPI: any) => void; // Define a type for coreAPI
}

// Placeholder for a module manager
// This would handle loading, registering, and managing modules.
// In a real scenario, this might fetch module bundles from a server or load them from a local directory.
export class ModuleManager {
    private modules: IPaneliumModule[] = [];

    registerModule(module: IPaneliumModule) {
        this.modules.push(module);
        if (module.register) {
            // module.register(this.getCoreAPI()); // Pass a core API object to the module
        }
        console.log(`Module "${module.name}" registered.`);
    }

    getModules() {
        return this.modules;
    }

    // Example: Get all routes from registered modules
    getModuleRoutes() {
        const allRoutes = {user: [], admin: []};
        this.modules.forEach(module => {
            if (module.routes?.user) {
                // @ts-ignore
                allRoutes.user.push(...module.routes.user);
            }
            if (module.routes?.admin) {
                // @ts-ignore
                allRoutes.admin.push(...module.routes.admin);
            }
        });
        return allRoutes;
    }

    // private getCoreAPI() {
    //   // Define and return an API object that modules can use to interact with the core system
    //   return {
    //     // Example: showNotification: (message: string) => { /* ... */ },
    //     // Example: getCurrentUser: () => { /* ... */ },
    //   };
    // }
}

export const moduleManager = new ModuleManager();

// Sample Module (Billing Example - conceptual)
const billingModule: IPaneliumModule = {
    name: "Billing System",
    version: "1.0.0",
    routes: {
        user: [
            // Example: { path: "billing", element: <BillingPage /> }
        ],
        admin: [
            // Example: { path: "billing-settings", element: <BillingAdminSettings /> }
        ]
    },
    components: {
        userDashboardWidgets: [/* Example: <BillingSummaryWidget /> */],
        adminSettingsTabs: [/* Example: <BillingSettingsTab /> */]
    },
    register: (coreAPI) => {
        console.log("Billing module registered and can use core API:", coreAPI);
        // coreAPI.showNotification("Billing module active!");
    }
};

// Sample Module (Proxmox VM Management Example - conceptual)
const proxmoxModule: IPaneliumModule = {
    name: "Proxmox Integration",
    version: "1.0.0",
    routes: {
        user: [
            // { path: "proxmox-vms", element: <ProxmoxUserPage /> }
        ],
        admin: [
            // { path: "proxmox-nodes", element: <ProxmoxAdminPage /> }
        ]
    },
    components: {
        // Could add a widget to the server details page to show Proxmox VM status
    },
    register: () => {
        console.log("Proxmox module registered.");
    }
};

// Register sample modules (in a real app, this would be dynamic)
// moduleManager.registerModule(billingModule);
// moduleManager.registerModule(proxmoxModule);

// How to potentially use module routes in your main routes.ts (conceptual):
// import { moduleManager } from './modules/shared';
// const moduleRoutes = moduleManager.getModuleRoutes();
// ... then spread moduleRoutes.user and moduleRoutes.admin into the respective layouts.

