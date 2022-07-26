// React imports
import React, { ChangeEvent } from 'react';

// Library imports
import { CdsControlMessage } from '@cds/react/forms';
import { FieldError } from 'react-hook-form';
import { CdsSelect } from '@cds/react/select';
import * as yup from 'yup';

interface ImageProps<T> {
    osImageTitle: string;
    field: string;
    errors: { [key: string]: FieldError | undefined };
    register: any;
    onOsImageSelected: (osImage: string, fieldName?: string) => void;
    images: T[];
}
export function osImagesValidation() {
    return yup.string().nullable().required('Please select an OS image');
}

function RetrieveOSImages<T>(props: ImageProps<T>) {
    const { osImageTitle, field, errors, register, images, onOsImageSelected } = props;
    const handleOsImageSelect = (event: ChangeEvent<HTMLSelectElement>) => {
        onOsImageSelected(event.target.value || '', field);
    };
    const fieldError = errors[field];
    return (
        <div cds-layout="m:lg">
            <h1>{osImageTitle}</h1>
            <CdsSelect layout="compact" controlWidth="shrink">
                <label>OS Image with Kubernetes </label>
                <select {...register(field)} onChange={handleOsImageSelect}>
                    {images.map((image) => (
                        <option key={image.name}>{image.name}</option>
                    ))}
                </select>
            </CdsSelect>
            {fieldError && (
                <div>
                    &nbsp;<CdsControlMessage status="error">{fieldError.message}</CdsControlMessage>
                </div>
            )}
        </div>
    );
}

export default RetrieveOSImages;
