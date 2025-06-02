import type {IPaneliumModule} from "../shared";

// This is a sample module to demonstrate the modular architecture.
// Modules can define their own routes, components, and register with the core system.

const sampleUserModule: IPaneliumModule = {
    name: "Sample User Module",
    version: "0.1.0",
    routes: {
        user: [
            // Path relative to the 'app' directory
            {path: "sample-feature", file: "modules/sample-module/SampleModulePage.tsx"},
        ],
    },
    components: {
        // Example: A widget to be displayed on the user dashboard
        // userDashboardWidgets: [() => <div>Sample Module Widget</div>],
    },
    register: (coreAPI) => {
        console.log("Sample User Module registered with core API:", coreAPI);
        // You could use coreAPI to interact with the main application,
        // for example, to add items to a menu, or subscribe to events.
    },
};

const sampleAdminModule: IPaneliumModule = {
    name: "Sample Admin Module",
    version: "0.1.0",
    routes: {
        admin: [
            // Path relative to the 'app' directory
            {path: "sample-admin-tool", file: "modules/sample-module/SampleModuleAdminPage.tsx"},
        ],
    },
    components: {
        // Example: A tab to be added to the admin settings page
        // adminSettingsTabs: [() => <div>Sample Module Admin Settings</div>],
    },
    register: (coreAPI) => {
        console.log("Sample Admin Module registered with core API:", coreAPI);
    },
};

// Register the sample modules
// In a real application, modules might be loaded dynamically (e.g., from separate bundles or plugins).
// moduleManager.registerModule(sampleUserModule);
// moduleManager.registerModule(sampleAdminModule);

// To make this module truly distributable, it would typically be its own package.
// For this example, we'll just export its definition.

export {sampleUserModule, sampleAdminModule};
