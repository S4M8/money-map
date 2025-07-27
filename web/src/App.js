import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { Plus, Trash2, DollarSign, TrendingUp, Calculator, Upload } from 'lucide-react';

const App = () => {
  const [income, setIncome] = useState([]);
  const [expenses, setExpenses] = useState([]);
  const [funds, setFunds] = useState({
    emergencyFund: 0.00,
    educationFund: 0.00,
    investments: 0.00,
    other: 0.00
  });

  const [newIncomeItem, setNewIncomeItem] = useState({ date: '', name: '', amount: '' });
  const [newExpenseItem, setNewExpenseItem] = useState({ date: '', name: '', amount: '', category: 'Core' });

  const [selectedMonth, setSelectedMonth] = useState(new Date().getMonth() + 1);
  const [selectedYear, setSelectedYear] = useState(new Date().getFullYear());

  const [isCsvModalOpen, setIsCsvModalOpen] = useState(false);
  const [csvFile, setCsvFile] = useState(null);

  const fetchData = () => {
    axios.get(`/api/income?month=${selectedMonth}&year=${selectedYear}`).then(res => setIncome(res.data || []));
    axios.get(`/api/expenses?month=${selectedMonth}&year=${selectedYear}`).then(res => setExpenses(res.data || []));
    axios.get('/api/funds').then(res => setFunds(res.data || {}));
  };

  useEffect(() => {
    fetchData();
  }, [selectedMonth, selectedYear]);

  // Calculations
  const totalIncome = income.reduce((sum, item) => sum + item.amount, 0);
  const totalExpenses = expenses.reduce((sum, item) => sum + item.amount, 0);
  const remainingAmount = totalIncome - totalExpenses;

  const coreExpenses = expenses.filter(e => e.category === 'Core').reduce((sum, item) => sum + item.amount, 0);
  const choiceExpenses = expenses.filter(e => e.category === 'Choice').reduce((sum, item) => sum + item.amount, 0);
  const compoundAmount = remainingAmount <= 0 ? 0 : remainingAmount;

  const corePercentage = totalIncome > 0 ? (coreExpenses / totalIncome) * 100 : 0;
  const choicePercentage = totalIncome > 0 ? (choiceExpenses / totalIncome) * 100 : 0;
  const compoundPercentage = totalIncome > 0 ? (compoundAmount / totalIncome) * 100 : 0;

  // Money Map Score calculation
  const coreScore = corePercentage <= 50 ? 1 : 0;
  const choiceScore = choicePercentage <= 30 ? 1 : 0;
  const compoundScore = compoundPercentage >= 20 ? 1 : 0;
  const totalScore = coreScore + choiceScore + compoundScore;

  const getScoreLabel = (score) => {
    switch(score) {
      case 3: return 'Great';
      case 2: return 'Okay';
      case 1: return 'Need Improvement';
      case 0: return 'Poor';
      default: return 'Poor';
    }
  };

  const getScoreColor = (score) => {
    switch(score) {
      case 3: return 'text-green-600 bg-green-100';
      case 2: return 'text-blue-600 bg-blue-100';
      case 1: return 'text-yellow-600 bg-yellow-100';
      case 0: return 'text-red-600 bg-red-100';
      default: return 'text-red-600 bg-red-100';
    }
  };

  // Auto-distribute compound funds
  useEffect(() => {
    if (compoundAmount > 0) {
      const perFund = compoundAmount / 3;
      const updatedFunds = {
        emergencyFund: perFund,
        educationFund: perFund,
        investments: perFund,
        other: 0
      };
      setFunds(updatedFunds);
      axios.put('/api/funds', updatedFunds);
    } else {
      const updatedFunds = {
        emergencyFund: 0,
        educationFund: 0,
        investments: 0,
        other: 0
      };
      setFunds(updatedFunds);
      axios.put('/api/funds', updatedFunds);
    }
  }, [compoundAmount]);

  const addIncomeItem = () => {
    if (newIncomeItem.name && newIncomeItem.amount && newIncomeItem.date) {
      axios.post('/api/income', {
        ...newIncomeItem,
        amount: parseFloat(newIncomeItem.amount) || 0
      }).then(() => {
        fetchData();
        setNewIncomeItem({ date: '', name: '', amount: '' });
      });
    }
  };

  const addExpenseItem = () => {
    if (newExpenseItem.name && newExpenseItem.amount && newExpenseItem.date) {
      axios.post('/api/expenses', {
        ...newExpenseItem,
        amount: parseFloat(newExpenseItem.amount) || 0
      }).then(() => {
        fetchData();
        setNewExpenseItem({ date: '', name: '', amount: '', category: 'Core' });
      });
    }
  };

  const removeIncomeItem = (id) => {
    axios.delete(`/api/income/${id}`).then(() => fetchData());
  };

  const removeExpenseItem = (id) => {
    axios.delete(`/api/expenses/${id}`).then(() => fetchData());
  };

  const handleCsvUpload = () => {
    if (csvFile) {
      const formData = new FormData();
      formData.append('file', csvFile);
      axios.post('/api/upload', formData).then(() => {
        fetchData();
        setIsCsvModalOpen(false);
      });
    }
  };

  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-6">
      <div className="max-w-7xl mx-auto">
        <div className="bg-white rounded-2xl shadow-xl p-8">
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold text-gray-800 mb-2">Monthly Accounting Dashboard</h1>
            <p className="text-gray-600">Track your income, expenses, and build wealth systematically</p>
          </div>

          {/* Summary Cards */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
            <div className="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-6 text-white">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-green-100 text-sm">Total Income</p>
                  <p className="text-2xl font-bold">{formatCurrency(totalIncome)}</p>
                </div>
                <DollarSign className="h-8 w-8 text-green-200" />
              </div>
            </div>

            <div className="bg-gradient-to-r from-red-500 to-red-600 rounded-xl p-6 text-white">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-red-100 text-sm">Total Expenses</p>
                  <p className="text-2xl font-bold">{formatCurrency(totalExpenses)}</p>
                </div>
                <TrendingUp className="h-8 w-8 text-red-200" />
              </div>
            </div>

            <div className="bg-gradient-to-r from-blue-500 to-blue-600 rounded-xl p-6 text-white">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-blue-100 text-sm">Remaining</p>
                  <p className="text-2xl font-bold">{formatCurrency(remainingAmount)}</p>
                </div>
                <Calculator className="h-8 w-8 text-blue-200" />
              </div>
            </div>

            <div className={`rounded-xl p-6 ${getScoreColor(totalScore)}`}>
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm opacity-75">Money Map Score</p>
                  <p className="text-2xl font-bold">{getScoreLabel(totalScore)}</p>
                  <p className="text-sm font-medium">{totalScore}/3</p>
                </div>
                <div className="text-2xl font-bold">{totalScore}</div>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Income Section */}
            <div className="bg-gray-50 rounded-xl p-6">
              <h2 className="text-2xl font-bold text-gray-800 mb-4 flex items-center">
                <DollarSign className="h-6 w-6 mr-2 text-green-600" />
                Capital (Income)
              </h2>
              
              <div className="space-y-3 mb-4">
                {income.map((item) => (
                  <div key={item.id} className="flex items-center justify-between bg-white rounded-lg p-3 shadow-sm">
                    <div className="flex-1">
                      <div className="font-medium text-gray-800">{item.name}</div>
                      <div className="text-sm text-gray-500">{new Date(item.date).toLocaleDateString()}</div>
                    </div>
                    <div className="flex items-center space-x-3">
                      <span className="font-bold text-green-600">{formatCurrency(item.amount)}</span>
                      <button
                        onClick={() => removeIncomeItem(item.id)}
                        className="text-red-500 hover:text-red-700 p-1"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                ))}
              </div>

              <div className="border-t pt-4">
                <div className="grid grid-cols-1 gap-2 mb-3">
                  <input
                    type="date"
                    value={newIncomeItem.date}
                    onChange={(e) => setNewIncomeItem({...newIncomeItem, date: e.target.value})}
                    className="border rounded-lg px-3 py-2"
                  />
                  <input
                    type="text"
                    placeholder="Income source name"
                    value={newIncomeItem.name}
                    onChange={(e) => setNewIncomeItem({...newIncomeItem, name: e.target.value})}
                    className="border rounded-lg px-3 py-2"
                  />
                  <input
                    type="number"
                    placeholder="Amount"
                    value={newIncomeItem.amount}
                    onChange={(e) => setNewIncomeItem({...newIncomeItem, amount: e.target.value})}
                    className="border rounded-lg px-3 py-2"
                  />
                </div>
                <button
                  onClick={addIncomeItem}
                  className="w-full bg-green-600 text-white rounded-lg py-2 hover:bg-green-700 flex items-center justify-center"
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Add Income
                </button>
              </div>
            </div>

            {/* Expenses Section */}
            <div className="bg-gray-50 rounded-xl p-6">
              <h2 className="text-2xl font-bold text-gray-800 mb-4">Core & Choice (Expenses)</h2>
              
              <div className="space-y-3 mb-4">
                {expenses.map((item) => (
                  <div key={item.id} className="flex items-center justify-between bg-white rounded-lg p-3 shadow-sm">
                    <div className="flex-1">
                      <div className="font-medium text-gray-800">{item.name}</div>
                      <div className="text-sm text-gray-500">{new Date(item.date).toLocaleDateString()}</div>
                      <span className={`text-xs px-2 py-1 rounded-full ${
                        item.category === 'Core' ? 'bg-blue-100 text-blue-800' : 'bg-purple-100 text-purple-800'
                      }`}>
                        {item.category}
                      </span>
                    </div>
                    <div className="flex items-center space-x-3">
                      <span className="font-bold text-red-600">{formatCurrency(item.amount)}</span>
                      <button
                        onClick={() => removeExpenseItem(item.id)}
                        className="text-red-500 hover:text-red-700 p-1"
                      >
                        <Trash2 className="h-4 w-4" />
                      </button>
                    </div>
                  </div>
                ))}
              </div>

              <div className="border-t pt-4">
                <div className="grid grid-cols-1 gap-2 mb-3">
                  <input
                    type="date"
                    value={newExpenseItem.date}
                    onChange={(e) => setNewExpenseItem({...newExpenseItem, date: e.target.value})}
                    className="border rounded-lg px-3 py-2"
                  />
                  <input
                    type="text"
                    placeholder="Expense name"
                    value={newExpenseItem.name}
                    onChange={(e) => setNewExpenseItem({...newExpenseItem, name: e.target.value})}
                    className="border rounded-lg px-3 py-2"
                  />
                  <input
                    type="number"
                    placeholder="Amount"
                    value={newExpenseItem.amount}
                    onChange={(e) => setNewExpenseItem({...newExpenseItem, amount: e.target.value})}
                    className="border rounded-lg px-3 py-2"
                  />
                  <select
                    value={newExpenseItem.category}
                    onChange={(e) => setNewExpenseItem({...newExpenseItem, category: e.target.value})}
                    className="border rounded-lg px-3 py-2"
                  >
                    <option value="Core">Core</option>
                    <option value="Choice">Choice</option>
                  </select>
                </div>
                <button
                  onClick={addExpenseItem}
                  className="w-full bg-red-600 text-white rounded-lg py-2 hover:bg-red-700 flex items-center justify-center"
                >
                  <Plus className="h-4 w-4 mr-2" />
                  Add Expense
                </button>
              </div>
            </div>
          </div>

          {/* Analysis Section */}
          <div className="mt-8 grid grid-cols-1 lg:grid-cols-2 gap-8">
            {/* Category Breakdown */}
            <div className="bg-gray-50 rounded-xl p-6">
              <h3 className="text-xl font-bold text-gray-800 mb-4">Category Analysis</h3>
              <div className="space-y-4">
                <div className="bg-white rounded-lg p-4">
                  <div className="flex justify-between items-center mb-2">
                    <span className="font-medium text-blue-600">Core Expenses</span>
                    <span className="font-bold">{formatCurrency(coreExpenses)}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div 
                      className="bg-blue-600 h-2 rounded-full" 
                      style={{width: `${Math.min(corePercentage, 100)}%`}}
                    ></div>
                  </div>
                  <div className="text-sm text-gray-600 mt-1">
                    {corePercentage.toFixed(1)}% of income {corePercentage <= 50 ? '✅' : '⚠️'}
                  </div>
                </div>

                <div className="bg-white rounded-lg p-4">
                  <div className="flex justify-between items-center mb-2">
                    <span className="font-medium text-purple-600">Choice Expenses</span>
                    <span className="font-bold">{formatCurrency(choiceExpenses)}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div 
                      className="bg-purple-600 h-2 rounded-full" 
                      style={{width: `${Math.min(choicePercentage, 100)}%`}}
                    ></div>
                  </div>
                  <div className="text-sm text-gray-600 mt-1">
                    {choicePercentage.toFixed(1)}% of income {choicePercentage <= 30 ? '✅' : '⚠️'}
                  </div>
                </div>

                <div className="bg-white rounded-lg p-4">
                  <div className="flex justify-between items-center mb-2">
                    <span className="font-medium text-green-600">Compound (Savings)</span>
                    <span className="font-bold">{formatCurrency(compoundAmount)}</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div 
                      className="bg-green-600 h-2 rounded-full" 
                      style={{width: `${Math.min(compoundPercentage, 100)}%`}}
                    ></div>
                  </div>
                  <div className="text-sm text-gray-600 mt-1">
                    {compoundPercentage.toFixed(1)}% of income {compoundPercentage >= 20 ? '✅' : '⚠️'}
                  </div>
                </div>
              </div>
            </div>

            {/* Fund Allocation */}
            <div className="bg-gray-50 rounded-xl p-6">
              <h3 className="text-xl font-bold text-gray-800 mb-4">Fund Allocation</h3>
              <div className="space-y-4">
                <div className="bg-white rounded-lg p-4">
                  <div className="flex justify-between items-center">
                    <span className="font-medium text-orange-600">Emergency Fund</span>
                    <span className="font-bold">{formatCurrency(funds.emergencyFund)}</span>
                  </div>
                </div>
                <div className="bg-white rounded-lg p-4">
                  <div className="flex justify-between items-center">
                    <span className="font-medium text-indigo-600">Education Fund</span>
                    <span className="font-bold">{formatCurrency(funds.educationFund)}</span>
                  </div>
                </div>
                <div className="bg-white rounded-lg p-4">
                  <div className="flex justify-between items-center">
                    <span className="font-medium text-green-600">Investments</span>
                    <span className="font-bold">{formatCurrency(funds.investments)}</span>
                  </div>
                </div>
                <div className="bg-white rounded-lg p-4">
                  <div className="flex justify-between items-center">
                    <span className="font-medium text-gray-600">Other</span>
                    <span className="font-bold">{formatCurrency(funds.other)}</span>
                  </div>
                </div>
                <div className="bg-gray-800 text-white rounded-lg p-4">
                  <div className="flex justify-between items-center">
                    <span className="font-medium">Total Funds</span>
                    <span className="font-bold">{formatCurrency(Object.values(funds).reduce((a, b) => a + b, 0))}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Floating Action Button */}
          <div className="fixed bottom-6 right-6">
            <button
              onClick={() => setIsCsvModalOpen(true)}
              className="bg-blue-600 text-white rounded-full p-4 shadow-lg hover:bg-blue-700"
            >
              <Upload className="h-6 w-6" />
            </button>
          </div>

          {/* CSV Upload Modal */}
          {isCsvModalOpen && (
            <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
              <div className="bg-white rounded-lg p-8">
                <h2 className="text-2xl font-bold mb-4">Upload CSV</h2>
                <input
                  type="file"
                  accept=".csv"
                  onChange={(e) => setCsvFile(e.target.files[0])}
                  className="mb-4"
                />
                <div className="flex justify-end">
                  <button
                    onClick={() => setIsCsvModalOpen(false)}
                    className="bg-gray-300 text-gray-800 rounded-lg px-4 py-2 mr-2"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={handleCsvUpload}
                    className="bg-blue-600 text-white rounded-lg px-4 py-2"
                  >
                    Upload
                  </button>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default App;
