import {create} from "zustand";
import {
    MatreshkaApi,
    ListConfigsRequest, ListConfigsResponse,
    InitReq
} from "@vervstack/matreshka";

export interface Api {
    key: string;
    pathPrefix: string;
    setKey: (key: string) => void;
    setPathPrefix: (pathPrefix: string) => void;

    initReq(): InitReq

    ListConfigs(req: ListConfigsRequest): Promise<ListConfigsResponse>;
}


export const useApi = create<Api>(
    (set, get) => ({
        key: setToLocalStorage(apiKeyKey(),
            loadFromLocalstorage(apiKeyKey()) || 'MatreshkaUI'),
        pathPrefix:
            setToLocalStorage(pathPrefixKey(),
                loadFromLocalstorage(pathPrefixKey()) || 'http://matreshka.vervstack.ru'),

        setKey(k: string) {
            setToLocalStorage(apiKeyKey(), k)
            set({key: k});
        },
        setPathPrefix(pp: string) {
            setToLocalStorage(pathPrefixKey(), pp)
            set({pathPrefix: pp})
        },

        ListConfigs(req: ListConfigsRequest): Promise<ListConfigsResponse> {
            const {initReq} = get();
            const ir = initReq();
            console.log(ir);

            return MatreshkaApi.ListConfigs(req, ir)
        },

        initReq(): InitReq {
            const {pathPrefix, key} = get();
            return {
                pathPrefix: pathPrefix,
                headers: {
                    'Grpc-Metadata-Authorization': key
                }
            }
        }
    }));


function pathPrefixKey() {
    return 'pathPrefix';
}

function apiKeyKey() {
    return 'apiKey';
}

function loadFromLocalstorage(pathPrefix: string) {
    return localStorage.getItem(pathPrefix);
}

function setToLocalStorage(key: string, val: string): string {
    localStorage.setItem(key, val);

    return val
}

