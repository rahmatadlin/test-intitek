# Toast Notifications System

This application includes a comprehensive toast notification system for providing real-time feedback to users.

## Overview

Toast notifications appear in the top-right corner of the screen and automatically disappear after 3 seconds (by default). They provide visual feedback for:

- ✅ Successful operations
- ❌ Error messages
- ℹ️ Informational messages
- ⚠️ Warning messages

## Implementation

### Composable (`useToast.js`)

The toast system is built using Vue 3's Composition API:

```javascript
import { useToast } from "../composables/useToast";

const toast = useToast();

// Success notification
toast.success("Operation completed successfully!");

// Error notification
toast.error("Something went wrong!");

// Info notification
toast.info("Here's some information");

// Warning notification
toast.warning("Please be careful!");
```

### Custom Duration

You can customize the duration (in milliseconds):

```javascript
// Show for 5 seconds
toast.success("This will show for 5 seconds", 5000);

// Show for 10 seconds
toast.error("This error will show for 10 seconds", 10000);
```

## Where Toasts Are Used

### Authentication

**Login (Login.vue)**

- ✅ Success: "Login successful! Welcome back."
- ❌ Error: Shows specific error message from API

**Logout (Layout.vue)**

- ℹ️ Info: "You have been logged out successfully."

### Product Management (Products.vue)

**Create Product**

- ✅ Success: "Product created successfully!"
- ❌ Error: Shows specific validation error

**Update Product**

- ✅ Success: "Product updated successfully!"
- ❌ Error: Shows specific validation error

**Delete Product**

- ✅ Success: "Product deleted successfully!"
- ❌ Error: "Failed to delete product"

**Export CSV**

- ✅ Success: "Products exported to CSV successfully!"
- ❌ Error: "Failed to export CSV"

**Generate Barcode**

- ❌ Error: "Failed to generate barcode"

**Load Products**

- ❌ Error: "Failed to load products"

### Dashboard (Dashboard.vue)

**Load Statistics**

- ❌ Error: "Failed to load dashboard statistics"

## Toast Types

### Success Toast (Green)

- **Background**: #10b981
- **Icon**: ✓
- **Use**: Confirm successful operations

### Error Toast (Red)

- **Background**: #ef4444
- **Icon**: ✕
- **Use**: Show errors and failures

### Info Toast (Blue)

- **Background**: #3b82f6
- **Icon**: ℹ
- **Use**: Provide informational messages

### Warning Toast (Orange)

- **Background**: #f59e0b
- **Icon**: ⚠
- **Use**: Show warnings

## Customization

### Styling

Toasts can be customized in `Toast.vue`:

```css
.toast-item {
  padding: 14px 16px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  /* ... */
}
```

### Animation

The toast uses slide-in animation from the right:

```css
@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}
```

## Features

1. **Auto-dismiss**: Toasts automatically disappear after specified duration
2. **Click to dismiss**: Users can click on a toast to dismiss it immediately
3. **Close button**: Each toast has a close (×) button
4. **Stacking**: Multiple toasts stack vertically
5. **Hover effect**: Toasts slightly move left on hover
6. **Responsive**: Works well on mobile devices

## Mobile Responsiveness

On mobile devices (< 768px), toasts:

- Span the full width with margins
- Adjust positioning for better visibility

```css
@media (max-width: 768px) {
  .toast-container {
    left: 10px;
    right: 10px;
    max-width: none;
  }
}
```

## Best Practices

1. **Keep messages concise**: Short, clear messages work best
2. **Use appropriate types**: Match the toast type to the message context
3. **Don't overuse**: Only show toasts for important feedback
4. **Provide context**: Include relevant details in error messages
5. **Test duration**: Ensure users have enough time to read messages

## Example Usage in New Components

```vue
<script setup>
import { useToast } from "../composables/useToast";

const toast = useToast();

const handleAction = async () => {
  try {
    // Perform action
    await someApiCall();
    toast.success("Action completed!");
  } catch (error) {
    toast.error(error.message || "Action failed");
  }
};
</script>
```

## Advanced Usage

### Multiple toasts

Toasts automatically stack:

```javascript
toast.success("First action completed");
toast.success("Second action completed");
toast.info("Here's some information");
```

### Programmatic removal

```javascript
const { toasts, removeToast } = useToast();

// Remove specific toast by ID
removeToast(toastId);
```

## Architecture

```
App.vue
  └── Toast.vue (Toast Container)
        └── useToast() composable
              ├── Global toast state
              ├── Toast management functions
              └── Toast array
```

The toast state is global and shared across all components using the composable pattern.
