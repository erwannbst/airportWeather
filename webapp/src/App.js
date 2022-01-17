import logo from './logo.svg';
import './App.css';
import React, {useState} from "react";
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";
import moment from 'moment';
import { LineChart, Line, CartesianGrid, XAxis, YAxis } from 'recharts';


function App() {
  const [selectValue, setSelectValue] = useState("Atmospheric pressure");
  const [selectDate, setSelectDate] = useState(new Date("2022-01-09"));
  const [startDate, setStartDate] = useState(new Date("2022-01-09"));
  const [endDate, setEndDate] = useState(new Date("2022-01-11"));
  const emptyData = () => {return <p>"No data"</p>};
  const [chart, setChart] = useState(emptyData);
  const [byDayData, setByDayData] = useState(emptyData);
  const [byTypeData, setByTypeData] = useState(emptyData);

  const statByDay = (date) => {
    setSelectDate(date);
    date = moment(date).format("yyyy-MM-DD");
    let url = "http://127.0.0.1:4000/airport_weather/average?day=" + date;
    fetch(url)
    .then(data => data.json())
    .then(json => {
      setByDayData(makeDataByDate(json, date));
    })
    console.log(date);
  }

  const makeDataByDate = (json, date) => {
      if(json.data){
      return (<>
        <table>
            <thead>
                <tr>
                    <th colSpan="2">{date}</th>
                </tr>
            </thead>
            <tbody>
              {json.data.data.map(row => {
               return <tr>
                          <td>{row.id}</td>
                          <td>{row.AvgValue}</td>
                      </tr>
              })}
            </tbody>
        </table>
        </>);
      }else{
        return <p>No data for this date</p>
      }
  }

  const measureByType = () => {
    let start = moment(startDate).format("yyyy-MM-DD");
    let end = moment(endDate).format("yyyy-MM-DD");
    console.log(selectValue);
    let type = encodeURI(selectValue);
    let url = `http://127.0.0.1:4000/airport_weather/measures_by_type?measure_type=${type}&start_date=${start}&end_date=${end}`
    console.log(url);
    fetch(url)
    .then(data => data.json())
    .then(json => {
      console.log(json);
      setChart(
        <LineChart
          width={900}
          height={600}
          data={json.data.data}
          margin={{
            top: 5,
            right: 30,
            left: 20,
            bottom: 5
          }}
        >
          <CartesianGrid strokeDasharray="3 3" />
          <XAxis dataKey="Time" />
          <YAxis />
          <Line type="monotone" dataKey="Value" stroke="#82ca9d" />
        </LineChart>
      )
    })
  }


  return (
    <div className="App">
      <header className="App-header">
        <p>
          Airport project
        </p>
      </header>
      <div id="body">
      <div className="part divG">
        <h1>Par jour</h1>
          <DatePicker selected={selectDate} onChange={date => statByDay(date)}/>
          {byDayData}
        </div>
      <div className="part divD">
        <h1>Measure by type</h1>
        <select value={selectValue} onChange={v => setSelectValue(v.target.value)}>
            <option value="Atmospheric pressure">Atmospheric pressure</option>
            <option value="Temperature">Temperature</option>
            <option value="Wind speed">Wind speed</option>
        </select>
        <DatePicker selected={startDate} onChange={date => setStartDate(date)}/>
        <DatePicker selected={endDate} onChange={date => setEndDate(date)}/>
        <button onClick={measureByType}>Afficher</button>
        {chart}
        </div>
        </div>
    </div>
  );
}

export default App;
