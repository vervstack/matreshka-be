import cls from '@/segments/Settings.module.css';

import cn from "classnames";

import {useApi} from "@/app/hooks/api/api.ts";
import Input from "@/components/shared/Input.tsx";
import {useEffect, useRef, useState} from "react";

export default function Settings() {
    const api = useApi()

    const [isHidden, setIsHidden] = useState(true);
    const containerRef = useRef<
        HTMLDivElement | null>(null);

    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            if (event.key === "Escape") {
                event.preventDefault();
                setIsHidden(true);
            }
        };

        const handleClickOutside = (event: MouseEvent) => {
            if (containerRef.current && containerRef.current.contains &&
                !containerRef.current.contains(event.target as Node)) {
                setIsHidden(true);
            }
        };

        document.addEventListener("keydown", handleKeyDown);
        document.addEventListener("mousedown", handleClickOutside);

        return () => {
            document.removeEventListener("keydown", handleKeyDown);
            document.removeEventListener("mousedown", handleClickOutside);
        };
    }, []);

    // Toggle visibility with Ctrl+. or Cmd+.
    useEffect(() => {
        const handleKeyDown = (event: KeyboardEvent) => {
            if ((event.ctrlKey || event.metaKey) && event.key === ".") {
                event.preventDefault();
                setIsHidden(prev => !prev);
            }
        };

        document.addEventListener("keydown", handleKeyDown);
        return () => {
            document.removeEventListener("keydown", handleKeyDown);
        };
    }, []);

    return (
        <div
            className={cn(cls.SettingsContainer, {
                [cls.hidden]: isHidden,
            })}
            ref={containerRef}
        >
            <Input
                label={'ApiKey'}
                inputValue={api.key}
                onChange={api.setKey}/>
            <Input
                label={'URL'}
                inputValue={api.pathPrefix}
                onChange={api.setPathPrefix}/>
        </div>
    )
}
