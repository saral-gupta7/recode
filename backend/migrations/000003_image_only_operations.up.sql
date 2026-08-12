ALTER TABLE jobs DROP CONSTRAINT jobs_operation_valid;

ALTER TABLE jobs ADD CONSTRAINT jobs_operation_valid CHECK (
    operation IN (
        'image_grayscale',
        'image_convert',
        'image_compress',
        'image_resize',
        'image_crop',
        'image_rotate',
        'image_flip',
        'image_thumbnail',
        'image_strip_metadata',
        'image_adjust',
        'image_blur',
        'image_sharpen',
        'image_pixelate',
        'image_padding'
    )
) NOT VALID;
