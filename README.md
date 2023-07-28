# Histogram and Binning

A histogram is a graphical representation of the distribution of a dataset. It is an approximate visualization of the distribution of numerical data. To construct a histogram, the range of values is divided into a series of intervals or "bins", and the frequency of data points that fall into each bin is then calculated.

However, when the data set has many unique values, the resulting histogram may become too large or complex. Binning is a technique used to simplify the histogram by reducing the number of bins. The process of binning involves merging adjacent bins, resulting in a histogram that is easier to analyze and visualize.

In this implementation, the binning process is triggered when the size of the histogram is greater than or equal to 50. When this occurs, the histogram is binned into 25, reducing the complexity of the histogram while still maintaining an overview of the data distribution.

## Binning Approach

The binning approach applied here is designed to be flexible and robust, able to handle a variety of data configurations:

1. **Handling Single Value Bins:** If a bin represents a single value rather than a range, this is taken into account during the binning process.

2. **Sorting Bins:** Before binning, the bins are sorted in ascending order to ensure that they are merged correctly.

3. **Dynamic Bin Size Calculation:** The size of the bins in the new histogram is determined based on the minimum and maximum values of the old histogram and the specified maximum number of bins. This ensures adaptability to different data ranges and bin sizes.

4. **Rounding to a Fixed Number of Significant Figures:** During the binning process, the range of each new bin is rounded to a consistent number of significant figures, simplifying the bin labels while maintaining an appropriate level of precision.

Please note that this rounding approach may result in non-uniform bin sizes as it prioritizes individual number precision over the precision of the difference between numbers. Additionally, the code does not handle edge cases where two different bins might round to the same value. In these cases, the bins would be merged, potentially reducing the number of bins below the target number. Handling these edge cases would require a more robust solution.

After the binning operation, users can get the total count of bins and the total count of all data points in the histogram, providing a summary of the data distribution.

### Building from Source
building the plugin requires all tags to be included

```bash
go build -gcflags="all=-N -l" -o histogram-debug.so -buildmode=plugin
```