# XLSX Processor Service

A Go-based service for processing Excel files with pagination and transformation capabilities.

## Overview

The XLSX Processor provides three main endpoints:

- **Paginate**: Splits multi-sheet Excel files into individual JSON files for each sheet
- **Transform**: Applies transformation rules to Excel files
- **TransformJson**: Converts JSON files to Excel format and applies transformation rules

## API Endpoints

### 1. Paginate Endpoint

**URL**: `POST /paginate`

Splits a multi-sheet Excel file into multiple single-sheet JSON files.

#### Request Body Structure

```json
{
  "input": {
    "storageType": "s3",
    "reference": {
      "id": "input-file-id",
      "bucket": "my-input-bucket",
      "prefix": "path/to/workbook.xlsx",
      "region": "us-east-1"
    },
    "credential": {
      "secrets": {
        "secret": "aws-secret-key",
        "accessToken": "aws-access-token"
      },
      "resources": {
        "id": "resource-id"
      }
    }
  },
  "output": {
    "storageType": "s3", 
    "reference": {
      "id": "output-file-id",
      "bucket": "my-output-bucket",
      "prefix": "output/paginated",
      "region": "us-east-1"
    },
    "credential": {
      "secrets": {
        "secret": "aws-secret-key",
        "accessToken": "aws-access-token"
      },
      "resources": {
        "id": "resource-id"
      }
    }
  }
}
```

#### Expected Response

```json
{
  "message": "Success: 3 sheets paginated",
  "attributes": {
    "sheetMinimals": [
      {
        "sheetName": "Sheet1",
        "sheetTabColor": "#FF0000"
      },
      {
        "sheetName": "Sheet2", 
        "sheetTabColor": "#00FF00"
      },
      {
        "sheetName": "Sheet3",
        "sheetTabColor": "#0000FF"
      }
    ],
    "textColors": ["#000000", "#FF0000"],
    "bgColors": ["#FFFFFF", "#FFFF00"]
  },
  "totalPages": 3
}
```

#### Output Files Created

For each sheet in the original Excel file, a separate JSON file will be created:

- `output/key/pages/1.json` - Contains Sheet1 data
- `output/`key `/pages/2.json` - Contains Sheet2 data
- `output/key/pages/3.json` - Contains Sheet3 data

Each JSON file contains a single sheet object:

```json
{
    "sheetName": "Sheet1",
    "cells": [
      [
        {"value": "Header1", "style": {...}},
        {"value": "Header2", "style": {...}}
      ],
      [
        {"value": "Data1", "style": {...}},
        {"value": "Data2", "style": {...}}
      ]
    ]
}
```

### 2. Transform Endpoint

**URL**: `POST /transform`

Applies transformation rules to Excel files (supports .xlsx files).

#### Request Body Structure

```json
{
  "input": {
    "storageType": "s3",
    "reference": {
      "id": "input-file-id", 
      "bucket": "my-input-bucket",
      "prefix": "path/to/workbook.xlsx",
      "region": "us-east-1"
    },
    "credential": {
      "secrets": {
        "secret": "aws-secret-key",
        "accessToken": "aws-access-token"
      },
      "resources": {
        "id": "resource-id"
      }
    }
  },
  "output": {
    "storageType": "s3",
    "reference": {
      "id": "output-file-id",
      "bucket": "my-output-bucket", 
      "prefix": "path/to/transformed.xlsx",
      "region": "us-east-1"
    },
    "credential": {
      "secrets": {
        "secret": "aws-secret-key",
        "accessToken": "aws-access-token"
      },
      "resources": {
        "id": "resource-id"
      }
    }
  },
  "rules": [
    {
      "pageCondition": {
        "sheetName": "Sheet1",
        "includeFormulas": false,
        "nonEmptyValueRedact": false
      },
      "actions": [
        {
          "operation": "value",
          "value": "sensitive",
          "actionType": "redact"
        },
        {
          "operation": "range", 
          "value": "C4:D9",
          "actionType": "redact"
        },
        {
          "operation": "column",
          "value": "E",
          "actionType": "exclude"
        }
      ]
    }
  ],
  "webhook": {
    "url": "https://my-app.com/webhook",
    "responseToken": "webhook-token",
    "payload": {
      "msg": "Transform completed",
      "browserTabID": "tab-123",
      "uuid": "process-uuid",
      "userId": "user-123",
      "s3Bucket": "my-output-bucket",
      "s3Key": "path/to/transformed.xlsx",
      "sourceId": "source-123",
      "status": "completed"
    }
  }
}
```

#### Available Rule Operations

- **`"value"`**: Redact cells containing specific values
- **`"range"`**: Redact cells in a specific range (e.g., "C4:D9")
- **`"textColor"`**: Redact cells with specific text color (hex color without #)
- **`"bgColor"`**: Redact cells with specific background color (hex color without #)
- **`"column"`**: Exclude entire columns (e.g., "C" or "E")
- **`"row"`**: Exclude entire rows (e.g., "4" or "10")

#### Expected Response

```json
{
  "message": "File transformed successfully"
}
```

### 3. TransformJson Endpoint

**URL**: `POST /transform-json`

Converts JSON files to Excel format and applies transformation rules. The input file must have a `.json` extension.

#### Request Body Structure

```json
{
  "input": {
    "storageType": "s3",
    "reference": {
      "id": "input-file-id",
      "bucket": "my-input-bucket", 
      "prefix": "path/to/data.json",
      "region": "us-east-1"
    },
    "credential": {
      "secrets": {
        "secret": "aws-secret-key",
        "accessToken": "aws-access-token"
      },
      "resources": {
        "id": "resource-id"
      }
    }
  },
  "output": {
    "storageType": "s3",
    "reference": {
      "id": "output-file-id",
      "bucket": "my-output-bucket",
      "prefix": "path/to/transformed.json",
      "region": "us-east-1"
    },
    "credential": {
      "secrets": {
        "secret": "aws-secret-key", 
        "accessToken": "aws-access-token"
      },
      "resources": {
        "id": "resource-id"
      }
    }
  },
  "rules": [
    {
      "pageCondition": {
        "sheetName": "Sheet1",
        "includeFormulas": true,
        "nonEmptyValueRedact": false
      },
      "actions": [
        {
          "operation": "value",
          "value": "confidential",
          "actionType": "redact"
        }
      ]
    }
  ]
}
```

#### Expected Response

```json
{
  "message": "File transformed successfully"
}
```

The output will be a JSON file containing the transformed sheet data in the same format as the paginate endpoint output.

## Request Body Field Descriptions

### Storage Types

Currently supported: `"s3"`

### Credentials

- **`secrets.secret`**: Storage service secret key
- **`secrets.accessToken`**: Storage service access token
- **`resources.id`**: Resource identifier

### Reference Object

- **`id`**: Unique identifier for the file reference
- **`bucket`**: Storage bucket name
- **`prefix`**: File path within the bucket
- **`region`**: Storage service region

### Page Condition

- **`sheetName`**: Target sheet name for rule application
- **`includeFormulas`**: Whether to include Excel formulas in processing
- **`nonEmptyValueRedact`**: Whether to redact all non-empty values

### Actions

- **`operation`**: Type of operation (value, range, textColor, bgColor, column, row)
- **`value`**: Target value/range/color for the operation
- **`actionType`**: Action to perform (currently supports "redact" and "exclude")

### Webhook (Optional)

Optional callback configuration for async processing notifications.

## Error Handling

The service returns appropriate HTTP status codes:

- **200/202**: Success
- **400**: Bad Request (validation errors, invalid file types)
- **500**: Internal Server Error (processing failures, storage errors)

Error responses include detailed error messages and may include rule/action indices for transformation errors.

## Limitations

- Maximum 1,000 sheets per Excel file for pagination
- JSON files must have `.json` extension for transform-json endpoint
- Color values in rules should be hex codes without the `#` prefix
