import cls from "@/segments/Loader.module.css";
import {useEffect, useState} from "react";
import {useToaster} from "@/app/hooks/toaster/Toaster.ts";

interface LoaderWrapperProps {
    children: React.JSX.Element[] | React.JSX.Element;

    load?: Promise<void>;
}

export default function LoaderWrapper({load, children}: LoaderWrapperProps) {
    const [isLoading, setIsLoading] = useState(false);
    const [err, setErr] = useState<Error | undefined>();
    const toaster = useToaster()

    useEffect(() => {
        if (err) {
            toaster.catchGrpc(err)
        }
    }, [err]);

    useEffect(() => {
        // TODO remove from eslint
        // eslint-disable-next-line react-hooks/set-state-in-effect
        setErr(undefined);

        if (!load) {
            setIsLoading(false);
            return
        }

        setIsLoading(true);
        load
            .catch(setErr)
            .finally(() => setIsLoading(false));
    }, [load]);

    return (
        <div className={cls.LoaderContainer}>
            {isLoading && <div className={cls.Spinner}/>}
            {err ? <div className={cls.Text}>{err.message}</div> : children}
        </div>
    )
}
