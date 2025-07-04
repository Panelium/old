class Page {
    #id;
    #component;

    constructor(id: string, component: React.FC) {
        this.#id = id;
        this.#component = component;
    }

    get id() {
        return this.#id;
    }

    get component() {
        return (id: string): React.FC => {
            if (id !== this.#id) {
                return () => (<></>);
            }

            const ThisComponent = this.#component;

            return () => (
                <div className="w-full h-170 p-4">
                    <ThisComponent/>
                </div>
            );
        }
    }
}

export default Page;