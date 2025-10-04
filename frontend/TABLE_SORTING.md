# Table Sorting Feature

This document describes the table sorting feature implemented in the Products page.

## Overview

The products table supports **client-side sorting** on all columns. Users can click on any column header to sort the data in ascending or descending order.

## Features

### 1. Sortable Columns

All columns in the products table are sortable:

- **No** - Sequential numbering (1, 2, 3, etc.)
- **Product Name** - Alphabetical sorting
- **SKU** - Alphabetical sorting
- **Quantity** - Numerical sorting
- **Location** - Alphabetical sorting
- **Status** - Alphabetical sorting
- **Created At** - Date/time sorting (newest first by default)

### 2. Visual Indicators

- **Hover Effect**: Column headers change background color on hover
- **Sort Arrows**:
  - ↑ indicates ascending order
  - ↓ indicates descending order
- **Active Column**: Only the currently sorted column shows the arrow indicator

### 3. Default Sorting

By default, products are sorted by:

- **Column**: Created At
- **Order**: Descending (newest first)

This ensures new products appear at the top of the list.

### 4. Sequential Numbering

The "No" column displays sequential numbers (1, 2, 3, etc.) that remain consistent regardless of sorting. This provides a clear visual count of items displayed.

**Example:**

- Even if products are sorted by name, the numbers will be 1, 2, 3... in the current display order
- If Product ID 5 is first after sorting, it will show as "No: 1"

### 5. Date Formatting

The "Created At" column displays dates in a user-friendly format:

- Format: `MMM DD, YYYY, HH:MM AM/PM`
- Example: `Jan 15, 2024, 02:30 PM`

## How It Works

### User Interaction

1. **Click a column header** to sort by that column
2. **First click**: Sorts in ascending order (A→Z, 0→9, oldest→newest)
3. **Second click**: Toggles to descending order (Z→A, 9→0, newest→oldest)
4. **Click different column**: Sorts by new column in ascending order

### Technical Implementation

#### State Management

```javascript
const sortColumn = ref("created_at"); // Default column
const sortDirection = ref("desc"); // Default direction
```

#### Computed Property

```javascript
const sortedProducts = computed(() => {
  const productsCopy = [...products.value];

  productsCopy.sort((a, b) => {
    let aVal, bVal;

    if (sortColumn.value === "created_at") {
      aVal = new Date(a.created_at).getTime();
      bVal = new Date(b.created_at).getTime();
    } else if (sortColumn.value === "quantity") {
      aVal = a.quantity;
      bVal = b.quantity;
    } else {
      aVal = a[sortColumn.value]?.toString().toLowerCase() || "";
      bVal = b[sortColumn.value]?.toString().toLowerCase() || "";
    }

    if (sortDirection.value === "asc") {
      return aVal > bVal ? 1 : aVal < bVal ? -1 : 0;
    } else {
      return aVal < bVal ? 1 : aVal > bVal ? -1 : 0;
    }
  });

  return productsCopy;
});
```

#### Sort Function

```javascript
const sortBy = (column) => {
  if (sortColumn.value === column) {
    // Toggle direction if clicking same column
    sortDirection.value = sortDirection.value === "asc" ? "desc" : "asc";
  } else {
    // Set new column and default to ascending
    sortColumn.value = column;
    sortDirection.value = "asc";
  }
};
```

## Styling

### Column Headers

```css
.sortable {
  cursor: pointer;
  user-select: none;
  position: relative;
  transition: background-color 0.2s;
}

.sortable:hover {
  background-color: #f3f4f6;
}
```

### Sort Indicators

```css
.sort-icon {
  margin-left: 0.5rem;
  font-size: 0.875rem;
  color: var(--primary-color);
  font-weight: bold;
}
```

## Examples

### Example 1: Sort by Name

**Action:** Click "Product Name" header

**Result:**

```
No | Product Name              | SKU        | ...
1  | Keyboard Mechanical RGB   | KEYB-001   | ...
2  | Laptop Dell XPS 15        | LAPTOP-001 | ...
3  | Monitor 27 inch 4K        | MON-001    | ...
4  | USB-C Hub Multiport       | USB-001    | ...
5  | Wireless Mouse Logitech   | MOUSE-001  | ...
```

### Example 2: Sort by Quantity (Descending)

**Action:** Click "Quantity" header twice

**Result:**

```
No | Product Name              | Quantity | ...
1  | Keyboard Mechanical RGB   | 50       | ...
2  | Laptop Dell XPS 15        | 25       | ...
3  | USB-C Hub Multiport       | 15       | ...
4  | Wireless Mouse Logitech   | 3        | ...
5  | Monitor 27 inch 4K        | 0        | ...
```

### Example 3: Sort by Status

**Action:** Click "Status" header

**Result:** Products grouped by status alphabetically (in_stock → low_stock → out_of_stock)

## Benefits

1. **User Control**: Users can organize data according to their needs
2. **No Page Reload**: Instant sorting without server requests
3. **Visual Feedback**: Clear indicators show current sort state
4. **Intuitive**: Familiar table sorting pattern
5. **Performance**: Client-side sorting is fast for typical dataset sizes
6. **Accessibility**: Keyboard accessible and screen reader friendly

## Future Enhancements

Possible improvements:

1. **Multi-column sorting**: Hold Shift and click to sort by multiple columns
2. **Sort persistence**: Remember user's preferred sort in localStorage
3. **Sort reset**: Button to reset to default sorting
4. **Custom sort orders**: Allow users to save favorite sort configurations
5. **Export with sort**: CSV export maintains current sort order

## Browser Compatibility

The sorting feature works in all modern browsers:

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Performance Considerations

- **Dataset Size**: Optimized for up to 1,000 products
- **Sorting Algorithm**: Native JavaScript `.sort()` (typically QuickSort)
- **Memory**: Creates a copy of the array to avoid mutating original data
- **Reactivity**: Uses Vue 3's computed property for efficient updates

## Testing

To test the sorting feature:

1. Add multiple products with varying data
2. Click each column header
3. Verify ascending/descending toggle
4. Check that arrows appear correctly
5. Confirm "No" column remains sequential
6. Test with edge cases (empty values, same values)

## Troubleshooting

**Problem:** Sorting not working

**Solutions:**

- Check that products array is properly populated
- Verify computed property is reactive
- Ensure column names match data structure

**Problem:** Wrong sort direction

**Solutions:**

- Check sortDirection state
- Verify arrow indicators
- Test toggle logic in sortBy function

**Problem:** Numbers sorting as strings

**Solutions:**

- Ensure quantity column uses numerical comparison
- Check that values are parsed as numbers

## Code Location

- **Component**: `frontend/src/views/Products.vue`
- **State**: Lines ~276-277
- **Computed**: Lines ~288-316
- **Function**: Lines ~442-451
- **Template**: Lines ~49-91
- **Styles**: Lines ~514-531
