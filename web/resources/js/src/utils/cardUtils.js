const getBgColorByProgress = (progress) => {
    const progressValue = parseInt(progress);
    if (progressValue < 25) return 'dark';
    if (progressValue >= 25 && progressValue < 50) return 'danger';
    if (progressValue >= 50 && progressValue < 75) return 'secondary';
    if (progressValue >= 75) return 'primary';
    return 'primary';
}

const getTextColorByProgress = (progress) => {
    const progressValue = parseInt(progress);
    if (progressValue < 25) return 'text-light';
    if (progressValue >= 25 && progressValue < 50) return 'text-light';
    if (progressValue >= 50 && progressValue < 75) return 'text-light';
    if (progressValue >= 75) return 'text-light';
    return 'text-dark';
}

// Simple memoization function
const memoize = (fn) => {
    const cache = new Map();
    return (...args) => {
        const key = JSON.stringify(args);
        if (cache.has(key)) {
            return cache.get(key);
        }
        const result = fn(...args);
        cache.set(key, result);
        return result;
    };
};

export const getCardColorsByProgress = memoize((progress) => {
    const bgColor = getBgColorByProgress(progress);
    const textColor = getTextColorByProgress(progress);
    return { bgColor, textColor };
});

