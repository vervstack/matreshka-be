import {useParams} from "react-router-dom";
import ConfigDisplayWidget from "@/widget/ConfigDisplayWidget.tsx";

export default function ConfigPage() {
    const {configName} = useParams();

    if (!configName) return null;

    return (
        <ConfigDisplayWidget configName={configName}/>
    );
}
